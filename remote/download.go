package remote

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type worker struct {
	Url       string
	File      *os.File
	Count     int64
	SyncWG    sync.WaitGroup
	TotalSize int64
	// Progress
}

func (r RemoteResource) downloadmulti(workerCount int64) error {
	file_sz, err := getSizeAndCheckRangeSupport(r.URL)
	if err != nil {
		return err
	}
	file_path := r.realpath()
	f, err := os.OpenFile(file_path, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	worker := worker{
		Url:       r.URL,
		File:      f,
		Count:     workerCount,
		TotalSize: file_sz,
	}

	var start, end int64
	var partial_size = int64(file_sz / workerCount)

	// now := time.Now().UTC()
	for num := int64(0); num < worker.Count; num++ {
		if num == worker.Count {
			end = file_sz // last part
		} else {
			end = start + partial_size
		}

		worker.SyncWG.Add(1)
		go worker.writeRange(num, start, end-1)
		start = end
	}

	worker.SyncWG.Wait()
	// log.Println("Elapsed time:", time.Since(now))
	return nil
}

func (w *worker) writeRange(part_num int64, start int64, end int64) {
	var written int64
	body, size, err := w.getRangeBody(start, end)
	if err != nil {
		log.Fatalf("Part %d request error: %s\n", part_num, err.Error())
	}
	defer body.Close()
	// defer w.Bars[part_num].Finish()
	defer w.SyncWG.Done()

	// Assign total size to progress bar
	// w.Bars[part_num].Total = size

	// New percentage flag
	percent_flag := map[int64]bool{}

	// make a buffer to keep chunks that are read
	buf := make([]byte, 4*1024)
	for {
		nr, er := body.Read(buf)
		if nr > 0 {
			nw, err := w.File.WriteAt(buf[0:nr], start)
			if err != nil {
				log.Fatalf("Part %d occured error: %s.\n", part_num, err.Error())
			}
			if nr != nw {
				log.Fatalf("Part %d occured error of short writiing.\n", part_num)
			}

			start = int64(nw) + start
			if nw > 0 {
				written += int64(nw)
			}

			// Update written bytes on progress bar
			// w.Bars[int(part_num)].Set64(written)

			// Update current percentage on progress bars
			p := int64(float32(written) / float32(size) * 100)
			_, flagged := percent_flag[p]
			if !flagged {
				percent_flag[p] = true
				// w.Bars[int(part_num)].Prefix(fmt.Sprintf("Part %d  %d%% ", part_num, p))
			}
		}
		if er != nil {
			if er.Error() == "EOF" {
				if size == written {
					// Download successfully
				} else {
					panic(fmt.Errorf("Part %d unfinished.\n", part_num))
				}
				break
			}
			panic(fmt.Errorf("Part %d occured error: %s\n", part_num, er.Error()))
		}
	}
}

func (w *worker) getRangeBody(start int64, end int64) (io.ReadCloser, int64, error) {
	var client http.Client
	req, err := http.NewRequest("GET", w.Url, nil)
	// req.Header.Set("cookie", "")
	// log.Printf("Request header: %s\n", req.Header)
	if err != nil {
		return nil, 0, err
	}

	// Set range header
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	size, err := strconv.ParseInt(resp.Header["Content-Length"][0], 10, 64)
	return resp.Body, size, err
}

func getSizeAndCheckRangeSupport(url string) (size int64, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	header := res.Header
	accept_ranges, supported := header["Accept-Ranges"]
	if !supported {
		return 0, errors.New("Doesn't support header `Accept-Ranges`.")
	} else if supported && accept_ranges[0] != "bytes" {
		return 0, errors.New("Support `Accept-Ranges`, but value is not `bytes`.")
	}
	size, err = strconv.ParseInt(header["Content-Length"][0], 10, 64)
	return
}
