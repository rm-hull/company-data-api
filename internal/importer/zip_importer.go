package importer

type ZipImporter interface {
	Import(zipPath string) error
}