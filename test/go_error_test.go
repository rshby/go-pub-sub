package test

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"sync"
	"testing"
	"time"
)

func TestErrorGoroutine(t *testing.T) {
	wg := &sync.WaitGroup{}

	errorChan := make(chan error, 3)
	//defer close(errorChan)
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go ErrorSpawning(wg, 2*time.Second, errorChan)
	}

	wg.Wait()
	close(errorChan)
	fmt.Println("channel diclose dulu")
	for err := range errorChan {
		if err != nil {
			logrus.Error(err)
		}
	}
}

func ErrorSpawning(wg *sync.WaitGroup, second time.Duration, errorChan chan error) {
	defer wg.Done()
	time.Sleep(second)
	errorChan <- errors.New("error from goroutine")
}

func TestWithoutBuffer(t *testing.T) {
	errorChan := make(chan error, 1)

	var counter int
	mu := &sync.Mutex{}
	go func() {
		for i := 0; i < 10; i++ {
			go func(i int, errorChan chan error) {
				mu.Lock()
				defer mu.Unlock()
				counter++
				if counter == 10 {
					errorChan <- errors.New("sudah terakhir")
					close(errorChan)
					return
				}

				errorChan <- errors.New("ini error")
			}(i, errorChan)
		}
	}()

	for err := range errorChan {
		logrus.Error(err)
	}
}

func TestGetTwoData(t *testing.T) {
	t.Run("test get 2 data paralel with error", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		errChan := make(chan error, 2)

		wg.Add(1)
		go func() {
			defer wg.Done()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
		}()

		wg.Wait()
		close(errChan)

		// cek error
		for err := range errChan {
			if err != nil {
				logrus.Error(err)
			}
		}
	})

	t.Run("test get 2 data paralel no error", func(t *testing.T) {
		type User struct {
			Id   int
			Name string
		}

		type Product struct {
			Id    int
			Name  string
			Price float64
		}

		// input params

		// get two data paralel
		wg := &sync.WaitGroup{}
		errorChan := make(chan error, 2)
		userChan := make(chan User, 1)
		productChan := make(chan Product, 1)

		wg.Add(1)
		go func() {
			defer wg.Done()

			// call from repository
			errorChan <- nil
			userChan <- User{1, "Reo"}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			// call from repository to get data product
			errorChan <- nil
			productChan <- Product{1, "iPhone", 22000000}
		}()

		wg.Wait()
		close(errorChan)
		//close(userChan)
		//close(productChan)

		for err := range errorChan {
			if err != nil {
				t.Fatalf("error : %v", err)
			}
		}

		// get
		user := <-userChan
		assert.NotNil(t, user)

		product := <-productChan
		assert.NotNil(t, product)
	})

	t.Run("test get 2 data assign from inside goroutine", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		errChan := make(chan error, 2)
		var user string
		var product string

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			user = "reo"
			errChan <- nil
		}(wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			product = "iphone"
			errChan <- nil
		}(wg)

		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				t.Fatalf("error : %v", err)
			}
		}

		// get data
		fmt.Println(user)
		fmt.Println(product)
	})

	t.Run("test get 2 data with error group", func(t *testing.T) {
		//wg := &sync.WaitGroup{}
		eg := &errgroup.Group{}

		eg.Go(func() error {
			return nil
		})

		eg.Go(func() error {
			return errors.New("kedua")
		})

		if err := eg.Wait(); err != nil {
			t.Fatalf("ada error : %v", err)
		}
	})
}

func TestChannelErrorBuffer(t *testing.T) {
	wg := &sync.WaitGroup{}
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- errors.New("error di sini")
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			logrus.Error(err)
		}
	}
}
