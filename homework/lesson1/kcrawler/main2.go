package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"golang.org/x/net/html"
)

// Объект который имеет все необходимые поля и методы для поиска
type VisitedPages struct {
	sync.Mutex
	wg *sync.WaitGroup
	// карта с url которые мы уже обработали
	visited map[string]string
	depthFromParent map[string]int
	maxDepth int
}

// Создаём новую сущность объекта
func NewVisitedPages(maxDepth int) *VisitedPages {
	return &VisitedPages{
		visited: make(map[string]string),
		depthFromParent: make(map[string]int),
		maxDepth: maxDepth,
		wg: new(sync.WaitGroup),
	}
}

// адрес в интернете
var url string
// насколько глубоко нам надо смотреть
var maxDepthProperty int
// Как вы помните, функция инициализации стартует первой
// Далее в курсе вы изучите ряд других возможностей чтения конфигурации,
// которые будут гораздо более удобные чем в данном примере
func init() {
	// задаём и парсим флаги
	flag.StringVar(&url, "url", "", "url address")
	flag.IntVar(&maxDepthProperty, "depth", 300, "max depth for analize")
	flag.Parse()
	// Проверяем обязательное условие
	if url == "" {
		log.Print("no url set by flag")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	started := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	go catchOs(cancel)
	// задаём количество ошибок достигнув которое работа приложения будет прервана
	maxErrorsBeforeChancel := 1113
	visitedPages := NewVisitedPages(maxDepthProperty)
	// создаём канал для результатов
	resultChan := make(chan string, 100)
	// создаём канал для ошибок
	errChan := make(chan error, 100)
	// создаём вспомогательные каналы для горутины которая будет вычитывать сообщения
	// из канала с результатами и с ошибками
	shutdownChanForReaders := make(chan struct{})
	readersDone := make(chan struct{})
	// запускаем горутину для чтения из каналов
	go startReaders(cancel, resultChan, errChan, shutdownChanForReaders, readersDone,
		maxErrorsBeforeChancel)
	// синхронизация окончания обхода анализа страниц через вэйтгруппу
	visitedPages.wg.Add(1)
	// запуск основной логики
	// внутри есть рекурсивные запуски анализа в других горутинах
	go visitedPages.analize(ctx, url, url, resultChan, errChan)
	// дожидаемся когда весь анализ окончен
	visitedPages.wg.Wait()
	// после окончания анализа мы могли успеть обработать не всю информацию из каналов
	// поэтому нам следует сообщить что новых данных не будет
	// конструкция с буферизированные каналами будет выглядеть намного проще,
	// но в нашем случае она не подходит
	shutdownChanForReaders <- struct{}{}
	// хороший тон всегда закрывать каналы, которые точно не будут использованы
	close(errChan)
	close(resultChan)
	// ждём завершения работы чтения в своей горутине
	<-readersDone
	log.Println(time.Since(started))
}
// ловим сигналы выключения
func catchOs(cancel context.CancelFunc) {
	osSignalChan := make(chan os.Signal)
	signal.Notify(osSignalChan, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-osSignalChan
	// если сигнал получен, отменяем контекст работы
	cancel()
}
func startReaders(cancel context.CancelFunc, resultChan chan string, errChan chan error,
	shutdownChanForReaders chan struct{}, readersDone chan struct{}, maxErrorsBeforeChancel int) {
	// начинаем цикл чтения из каналов
	var errCount int
	for {
		// порядок внутри выбора важен. Сообщение что пора выключаться придёт только после
		// того как другие каналы будут пустыми
		select {
		case result := <-resultChan:
			log.Printf("crawling result: %v", result)
		case err := <-errChan:
			log.Printf("when crawling got error: %v", err)
			errCount++
			if errCount == maxErrorsBeforeChancel {
				cancel()
			}
			// если мы дошли до этой части, значит пора прекращать работу из-за количества
			// ошибок
			// заметьте что мы не обратабывам прерывание контекста в уже полученных
			// данных, т.е. не отбрасываем полученную информацию
		case <-shutdownChanForReaders:
			// отправляем сигнал чтение из каналов прекращено
			readersDone <- struct{}{}
			return
		}
	}
}
// рекурсивно сканируем страницы

func (visitedPages *VisitedPages) analize(ctx context.Context, url, baseurl string, resultChan chan string,
	errChan chan error) {
	// При использовании синхронизаций при помощи каналов или
	// вэйтгруп лучше как можно раньше объявлять отложенные вызовы (defer)
	// Это может не раз спасти вас, ведь даже когда где-то внутри неожиданно
	// возникнет паника по - мы всё равно сможем отработать более-менее корректно
	defer visitedPages.wg.Done()
	// проверяем что контекст исполнения актуален
	select {
	case <-ctx.Done():
		errChan <- fmt.Errorf("cancel analize page %s", url)
		return
	default:
		// проверка глубины
		if visitedPages.isMaxDepth() {
			return
		}
		page, err := parse(url)
		if err != nil {
			// ошибку отправляем в канал ошибок, а не обрабатываем на месте
			errChan <- fmt.Errorf("error when getting page %s: %s", url, err)
			return
		}
		title := pageTitle(page)
		links := pageLinks(nil, page)
		// блокировка требуется, т.к. мы модифицируем карту объекта в рекурсии с новыми
		// гоуртинами
		visitedPages.Lock()
		if _, inMap := visitedPages.visited[url]; !inMap {
			// проверяем уникальность url
			visitedPages.visited[url] = title
			// отправляем результат в канал, не обрабатывая на месте
			resultChan <- fmt.Sprintf("%s -> %s\n", url, title)
		}
		visitedPages.Unlock()
		// рекурсивно ищем ссылки
		for link := range links {
			visitedPages.checkAndRecurseCall(ctx, link, baseurl, resultChan, errChan)
		}
	}
}
func (visitedPages *VisitedPages) checkAndRecurseCall(ctx context.Context, newURL, baseurl string,
	resultChan chan string, errChan chan error) {
	visitedPages.Lock()
	defer visitedPages.Unlock()
	// если ссылка найдена, то запускаем анализ по новой ссылке
	if visitedPages.visited[newURL] == "" && strings.HasPrefix(newURL, baseurl) {
		visitedPages.wg.Add(1)
		go visitedPages.analize(ctx, newURL, baseurl, resultChan, errChan)
	}
}
// парсим страницу
func parse(url string) (*html.Node, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't get page")
	}
	b, err := html.Parse(r.Body)
	if err != nil {
		return nil, fmt.Errorf("can't parse page")
	}
	return b, err
}
// ищем заголовок на странице
func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}
// ищем все ссылки на страницы. Используем карту чтобы избежать дубликатов
func pageLinks(links map[string]struct{}, n *html.Node) map[string]struct{} {
	if links == nil {
		links = make(map[string]struct{})
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if _, inMap := links[a.Val]; !inMap {
					links[a.Val] = struct{}{}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = pageLinks(links, c)
	}
	return links
}
// проверяем не слишком ли глубоко мы нырнули
func (visitedPages *VisitedPages) isMaxDepth() bool {
	visitedPages.Lock()
	defer visitedPages.Unlock()
	if visitedPages.maxDepth <= visitedPages.depthFromParent[url] {
		return true
	}
	visitedPages.depthFromParent[url]++
	return false
}
