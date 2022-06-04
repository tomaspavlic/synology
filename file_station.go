package synology

const (
	fileStationListApiName        = "SYNO.FileStation.List"
	fileStationApiVersion         = 1
	fileStationApiListShareMethod = "list_share"
)

type FileStation struct {
	core *SynologyCore
}

type FileShare struct {
	IsDir bool
	Name  string
	Path  string
}

type FileShareListResponse struct {
	Offset int
	Total  int
	Shares []FileShare
}

func (s *FileStation) ListShares() ([]FileShare, error) {
	response, err := s.core.makeRequest(
		entryPath,
		fileStationListApiName,
		fileStationApiListShareMethod,
		fileStationApiVersion,
		nil,
	)

	if err != nil {
		return nil, err
	}
	data, err := unmarshal[FileShareListResponse](response)

	return data.Shares, err
}
