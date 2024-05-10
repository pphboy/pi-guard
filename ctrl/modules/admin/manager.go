package admin

// login 用jwt
// create
// list
// delete
import (
	"crypto/sha256"
	"fmt"
	"go-ctrl/db"
	"go-ctrl/models"
	"go-ctrl/modules/jwt"
	"go-node/tool"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Manager struct {
	db     *gorm.DB
	jjwwtt *jwt.MyG[models.PiManager]
}

func NewManager() *Manager {

	m := &Manager{
		db:     db.DB(),
		jjwwtt: jwt.NewMyG[models.PiManager](),
	}

	m.init()

	return m
}
func (s *Manager) init() {
	l, err := s.List()
	if err != nil {
		logrus.Fatal(err)
	}

	for _, v := range l {
		if strings.Compare("ADMIN", strings.ToUpper(v.ManagerUsername)) == 0 {
			return
		}
	}

	if err := s.Create(&models.PiManager{
		ManagerType:     models.USER_ADMIN,
		ManagerUsername: "ADMIN",
		ManagerPswd:     "admin",
	}); err != nil {
		logrus.Fatal("init user admin,", err)
	}

}

func (s *Manager) Login(t *models.PiManager) (token string, err error) {
	logrus.Printf("login manager %+v",t) 
	
	spm := models.PiManager{}
	if err := s.db.Unscoped().Where("MANAGER_USERNAME = ?", strings.ToUpper(t.ManagerUsername)).First(&spm).Error; err != nil {
		return "", err
	}

	logrus.Debug("db pswd", spm.ManagerPswd)
	neoPaswd := fmt.Sprintf("%x",sha256.Sum256([]byte(t.ManagerPswd+spm.ManagerSolt)))
	logrus.Debug("neo Password",neoPaswd)
	if strings.Compare(neoPaswd, spm.ManagerPswd) == 0 {
		return s.jjwwtt.GetJwtToken(t.ManagerUsername, spm), nil
	} else {
		return "", fmt.Errorf("fail password or username ")
	}
}

func (s *Manager) ValitedToken(token string) bool {
	g, err := s.jjwwtt.ParseToken(token)
	if err != nil {
		return false
	}

	_, ok := g["ManagerUsername"]
	return ok
}

func (s *Manager) Create(t *models.PiManager) error {
	t.ManagerId = tool.GetUUIDUpper()
	t.ManagerSolt = tool.GetUUIDUpper()
	t.ManagerPswd = fmt.Sprintf("%x", sha256.Sum256([]byte(t.ManagerPswd+t.ManagerSolt)))
	t.ManagerUsername = strings.ToUpper(t.ManagerUsername)
	logrus.Printf("create manager %+v",t) 

	if err := s.db.FirstOrCreate(t).Error; err != nil {
		return err
	}
	return nil
}

func (s *Manager) Update(t *models.PiManager) error {

	t.ManagerSolt = tool.GetUUIDUpper()
	t.ManagerPswd = fmt.Sprintf("%x", sha256.Sum256([]byte(t.ManagerPswd+t.ManagerSolt)))
	t.ManagerUsername = strings.ToUpper(t.ManagerUsername)

	if err := s.db.Model(t).
		Where("manager_id = ?", t.ManagerId).
		Where("MANAGER_USERNAME = ?", t.ManagerUsername).
		Updates(t).Error; err != nil {
		return err
	}
	return nil
}

func (s *Manager) Get(id string) (*models.PiManager, error) {
	pm := &models.PiManager{}
	if err := s.db.Where("manager_id = ?", id).First(pm).Error; err != nil {
		return nil, err
	}
	return pm, nil
}

func (s *Manager) Delete(id string) error {
	m, err := s.Get(id)
	if err != nil {
		return err
	}
	if strings.Compare("ADMIN", strings.ToUpper(m.ManagerUsername)) == 0 {
		return fmt.Errorf("fuck off, cant remove admin")
	}

	if err := s.db.Unscoped().Where("manager_id = ?", id).Delete(&models.PiManager{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *Manager) List() (ts []*models.PiManager, err error) {
	if err := s.db.Model(&models.PiManager{}).Find(&ts).Error; err != nil {
		return nil, err
	}
	return
}
