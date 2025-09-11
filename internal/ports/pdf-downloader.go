package ports

type PDFDownloader interface {
	Download(url, savePath string) error
}
