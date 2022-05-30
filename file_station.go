package synology

const (
	fileStationApiName = "SYNO.FileStation.List"
	fileStationVersion = 1
	listShareMethod    = "list_share"
)

type FileStation struct {
	core    *SynologyCore
	listApi *Api
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
	response, err := s.core.makeRequest(s.listApi.Path, s.listApi.Name, listShareMethod, fileStationVersion, nil)
	if err != nil {
		return nil, err
	}
	data, err := unmarshal[FileShareListResponse](response)

	return data.Shares, err
}
