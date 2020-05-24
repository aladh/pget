package download

import (
	"log"
	"os"
	"testing"
)

func BenchmarkDownload(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := Run("http://ipv4.download.thinkbroadband.com/1GB.zip", 50, false)
		if err != nil {
			b.Fatalf("error downloading file: %s", err)
		}

		b.StopTimer()

		err = os.Remove("1GB.zip")
		if err != nil {
			log.Printf("error removing file: %s\n", err)
		}

		b.StartTimer()
	}
}
