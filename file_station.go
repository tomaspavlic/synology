package main

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

func (s *SynologyCore) ListShares() ([]FileShare, error) {
	api, err := s.Find("SYNO.FileStation.List")
	if err != nil {
		return nil, err
	}

	response, err := s.makeRequest(api.Path, api.Name, "list_share", 1, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := readResponse[FileShareListResponse](response)

	return data.Shares, err
}
