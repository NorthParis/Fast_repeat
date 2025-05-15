//fetchall параллельная выбора URL с сообщением
//о затраченном времени и размере ответа для каждого из URL

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()     // засекаем время
	ch := make(chan string) // создаем канал
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // запуск горутины
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

}

// создаем подпрограмму горутины с get запросом и передачей в канал
func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // отправляем в канал ch
		return
	}
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close() // исключение утечки ресурсов
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
