package files

type fileManager interface {
	GetFilesInPath(path string) ([]string, error)
	GetFilesRecursivelyInPath(path string) ([]string, error)
	DoesFileExist(path string) (bool, error)
	DoesPathExist(path string) (bool, error)
	MoveFile(sourcePath, destinationPath string) error
	CopyFile(sourcePath, destinationPath string) error
}

type Service struct {
	manager fileManager
}

func NewService(manager fileManager) *Service {
	return &Service{
		manager: manager,
	}
}

func (s *Service) GetFilesInPath(path string) ([]string, error) {
	return s.manager.GetFilesInPath(path)
}

func (s *Service) GetFilesRecursivelyInPath(path string) ([]string, error) {
	return s.manager.GetFilesRecursivelyInPath(path)
}

func (s *Service) DoesFileExist(path string) (bool, error) {
	return s.manager.DoesFileExist(path)
}

func (s *Service) DoesPathExist(path string) (bool, error) {
	return s.manager.DoesPathExist(path)
}

func (s *Service) MoveFile(sourcePath, destinationPath string) error {
	return s.manager.MoveFile(sourcePath, destinationPath)
}

func (s *Service) CopyFile(sourcePath, destinationPath string) error {
	return s.manager.CopyFile(sourcePath, destinationPath)
}
