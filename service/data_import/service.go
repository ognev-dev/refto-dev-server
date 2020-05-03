package dataimport

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/ognev-dev/bits/config"
	"github.com/ognev-dev/bits/database"
	"github.com/ognev-dev/bits/database/model"
	"github.com/ognev-dev/bits/service/entity"
	entitytopic "github.com/ognev-dev/bits/service/entity_topic"
	"github.com/ognev-dev/bits/service/topic"
)

const (
	dataExt   = ".yaml"
	sampleExt = ".sample.yaml"
)

func Process() (err error) {
	// Mark all data as deleted,
	// and while importing restore existing entities
	// Entities still marked as deleted after import should be deleted for real
	_, err = database.ORM().Exec("UPDATE entities SET deleted_at=NOW()")
	if err != nil {
		return
	}

	// delete topics
	_, err = database.ORM().Exec("DELETE FROM entity_topics")
	if err != nil {
		return
	}
	_, err = database.ORM().Exec("DELETE FROM topics")
	if err != nil {
		return
	}

	err = importEntitiesFromDir()
	if err != nil {
		return
	}

	_, err = database.ORM().Exec("DELETE FROM entities WHERE deleted_at IS NOT NULL")
	if err != nil {
		return
	}

	return
}

func importEntitiesFromDir() (err error) {
	err = filepath.Walk(config.Get().Dir.Data, func(path string, f os.FileInfo, wErr error) (err error) {
		if wErr != nil {
			return wErr
		}

		// skip dirs, samples and non-yaml files
		if f.IsDir() || strings.HasSuffix(f.Name(), sampleExt) || !strings.HasSuffix(f.Name(), dataExt) {
			return
		}

		nameParts := strings.Split(f.Name(), ".")
		var eType string
		if len(nameParts) > 2 {
			eType = nameParts[len(nameParts)-2]
		}

		fData, err := ioutil.ReadFile(path)
		if err != nil {
			return
		}

		data := Data{}
		err = yaml.Unmarshal(fData, &data)
		if err != nil {
			return
		}

		extWithType := eType + dataExt
		if extWithType[0] != '.' {
			extWithType = "." + extWithType
		}
		token := strings.TrimPrefix(path, config.Get().Dir.Data)
		token = strings.TrimSuffix(token, extWithType)

		jsonData, err := yaml.YAMLToJSON(fData)
		if err != nil {
			return
		}

		entityEl := model.Entity{
			Token:     token,
			Title:     data.Title,
			Type:      eType,
			Data:      string(jsonData),
			CreatedAt: time.Now(),
		}
		err = entity.CreateOrUpdate(&entityEl)
		if err != nil {
			return
		}

		for _, name := range data.Topics {
			topicEl, err := topic.FirstOrCreate(name)
			if err != nil {
				return err
			}
			err = entitytopic.Create(model.EntityTopic{
				EntityID: entityEl.ID,
				TopicID:  topicEl.ID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return
}
