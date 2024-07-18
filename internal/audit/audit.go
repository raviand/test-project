package audit

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
)

type Auditory interface {
	Run()
}

type Auditor struct {
	auditChannel chan AuditLog
	wg           *sync.WaitGroup
}

type AuditLog struct {
	Method    string
	Path      string
	TimeStamp string
	User      string
}

func NewAuditRoutine(wg *sync.WaitGroup) (chan AuditLog, Auditory) {
	auditChannel := make(chan AuditLog, 5)

	return auditChannel, &Auditor{
		auditChannel: auditChannel,
		wg:           wg,
	}
}

func (a *Auditor) Run() {
	defer a.wg.Done()

	file, err := os.OpenFile("audit.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open audit file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	if fileInfo.Size() == 0 {
		headers := []string{"Method", "Path", "TimeStamp", "User"}
		if err := writer.Write(headers); err != nil {
			log.Printf("Failed to write headers to CSV: %v", err)
		}
	}

	for audit := range a.auditChannel {
		record := []string{audit.Method, audit.Path, audit.TimeStamp, audit.User}
		if err := writer.Write(record); err != nil {
			log.Printf("Failed to write record to CSV: %v", err)
		}
		writer.Flush() // Asegurar que los datos se escriban inmediatamente
		if err := writer.Error(); err != nil {
			log.Printf("Error flushing to CSV: %v", err)
		}
	}
}
