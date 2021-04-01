package main

type res struct {
	hash string
	err error
}

func myFilesHasher(done <-chan interface{}, files...string) <-chan res {
	resCh := make(chan res)

	go func() {
		defer close(resCh)
		for _, f := range files {
			hash, err := compute(f)
			result := res {
				hash: hash,
				err: err,
			}

			select {
			case <- done:
				return
				case resCh<-result:
			}
		}
	}()

	return resCh
}