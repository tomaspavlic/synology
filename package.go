package synology

const (
	// SYNO.Core.Package
	packageApiName       = "SYNO.Core.Package"
	packageApiVersion    = 2
	packageApiListMethod = "list"

	// SYNO.Core.Package.Control
	packageControlApiVersion     = 1
	packageControlApiName        = "SYNO.Core.Package.Control"
	packageControlApiStartMethod = "start"
	packateControlApiStopMethod  = "stop"
)

type PackageControl struct {
	*SynologyCore
}

type Package struct {
	core *SynologyCore

	Id      string
	Name    string
	Version string
}

type packageListResponse struct {
	Packages []*Package
}

type packageControlResponse struct {
	Message       string
	WorkerMessage []string
}

func (s *Package) changeState(method string) error {
	response, err := s.core.makeRequest(
		entryPath,
		packageControlApiName,
		method,
		packageControlApiVersion,
		parameters{"id": s.Id})

	if err != nil {
		return err
	}

	_, err = unmarshal[packageControlResponse](response)

	return err
}

func (s *Package) Stop() error {
	return s.changeState(packateControlApiStopMethod)
}

func (s *Package) Start() error {
	return s.changeState(packageControlApiStartMethod)
}

func (s *PackageControl) List() ([]*Package, error) {
	response, err := s.SynologyCore.makeRequest(
		entryPath,
		packageApiName,
		packageApiListMethod,
		packageApiVersion,
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := unmarshal[packageListResponse](response)

	// set api core reference to each package object
	for _, p := range resp.Packages {
		p.core = s.SynologyCore
	}

	return resp.Packages, err
}
