package defs

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"time"
)

// BytesCounter implements io.Reader and io.Writer interface, for counting bytes being read/written in HTTP requests
type BytesCounter struct {
	start   time.Time
	total   int
	payload []byte
	reader  io.ReadSeeker
}

// Write implements io.Writer
func (c *BytesCounter) Write(p []byte) (int, error) {
	n := len(p)
	c.total += n

	return n, nil
}

// Read implements io.Reader
func (c *BytesCounter) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	c.total += n

	return n, err
}

// Average returns the average bytes/second
func (c *BytesCounter) Average() float64 {
	return float64(c.total) / time.Now().Sub(c.start).Seconds()
}

func (c *BytesCounter) AvgMbits() string {
	return fmt.Sprintf("%.02f Mbps", c.Average()/131072)
}

func (c *BytesCounter) AvgHumanize() string {
	val := c.Average()

	if val < 1024 {
		return fmt.Sprintf("%.2f bytes/s", val)
	} else if val/1024 < 1024 {
		return fmt.Sprintf("%.2f KB/s", val/1024)
	} else if val/1024/1024 < 1024 {
		return fmt.Sprintf("%.2f MB/s", val/1024/1024)
	} else {
		return fmt.Sprintf("%.2f GB/s", val/1024/1024/1024)
	}
}

// GenerateBlob generates a random byte array of `uploadSize` in the `payload` field, and sets the `reader` field to
// read from it
func (c *BytesCounter) GenerateBlob() {
	c.payload = getRandomData(uploadSize)
	c.reader = bytes.NewReader(c.payload)
}

// ResetReader resets the `reader` field to 0 position
func (c *BytesCounter) ResetReader() (int64, error) {
	return c.reader.Seek(0, 0)
}

// Start will set the `start` field to current time
func (c *BytesCounter) Start() {
	c.start = time.Now()
}

// Total returns the total bytes read/written
func (c *BytesCounter) Total() int {
	return c.total
}

// CurrentSpeed returns the current bytes/second
func (c *BytesCounter) CurrentSpeed() float64 {
	return float64(c.total) / time.Now().Sub(c.start).Seconds()
}

// SeekWrapper is a wrapper around io.Reader to give it a noop io.Seeker interface
type SeekWrapper struct {
	io.Reader
}

// Seek implements the io.Seeker interface
func (r *SeekWrapper) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

// getAvg returns the average value of an float64 array
func getAvg(vals []float64) float64 {
	var total float64
	for _, v := range vals {
		total += v
	}

	return total / float64(len(vals))
}

// getRandomData returns an `length` sized array of random bytes
func getRandomData(length int) []byte {
	data := make([]byte, length)
	if _, err := rand.Read(data); err != nil {
		log.Fatalf("Failed to generate random data: %s", err)
	}
	return data
}