COMMON := observer.go worker.go

all: observerGO sanitizeGO speedGO

observerGO: 
	go build -o observerGO ./observer

sanitizeGO:
	go build -o sanitizeGO ./sanitize

speedGO:
	go build -o speedGO ./speed

clean:
	rm -f observerGO sanitizeGO speedGO