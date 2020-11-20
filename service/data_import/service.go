package dataimport

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/refto/server/config"
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
	"github.com/refto/server/service/data"
	"github.com/refto/server/service/entity"
	entitytopic "github.com/refto/server/service/entity_topic"
	"github.com/refto/server/service/topic"
)

// type of data is added as topic
// but some exceptions exist
var autoTopicExceptions = []string{
	// software is too broad topic
	// each data should expand meaning of "software" in topics
	// for example: database, chat, media-player, etc
	// and topic "software" becomes redundant
	"software",

	// definition is special kind of data that describes topic
	// and returned if only one topic selected
	// so "definition" must be hidden, as it is makes no sense
	"definition",
}

// TODO partial import
func Import() (err error) {
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

	err = importDataFromDir()
	if err != nil {
		return
	}

	// TODO: since collections introduced it not good to hard delete entities,
	//  because that will be confusing when entity from collection simply disappear.
	//  Maybe if entity is in any collection then keep it marked as deleted and display it as deleted when viewed in collection
	//  Or implement notification system and inform everyone (that have deleted entity in their collection) about deleted entity
	_, err = database.ORM().Exec("DELETE FROM entities WHERE deleted_at IS NOT NULL")
	if err != nil {
		return
	}

	return
}

func importDataFromDir() (err error) {
	err = filepath.Walk(config.Get().Dir.Data, func(path string, f os.FileInfo, wErr error) (err error) {
		if wErr != nil {
			return wErr
		}

		// skip dirs, samples, schemas and non-yaml files
		if !data.IsDataFile(f) {
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

		dataEl := Data{}
		err = yaml.Unmarshal(fData, &dataEl)
		if err != nil {
			return
		}

		extWithType := eType + data.YAMLExt
		if extWithType[0] != '.' {
			extWithType = "." + extWithType
		}
		token := strings.TrimPrefix(path, config.Get().Dir.Data)
		token = strings.TrimSuffix(token, extWithType)

		var entityData model.EntityData
		err = yaml.Unmarshal(fData, &entityData)
		if err != nil {
			return
		}

		// prepend
		prependTopics := make([]string, 0)
		if eType != "" && useTypeAsTopic(eType) {
			prependTopics = append(prependTopics, eType)
		}
		if eType == "definition" {
			prependTopics = append(prependTopics, nameParts[len(nameParts)-3])
		}

		if len(prependTopics) > 0 {
			entityData = addTypeToTopics(entityData, prependTopics...)
			dataEl.Topics = append([]string{eType}, dataEl.Topics...)
		}

		entityEl := model.Entity{
			Token:     token,
			Title:     dataEl.Title,
			Type:      eType,
			Data:      entityData,
			CreatedAt: time.Now(),
		}
		err = entity.CreateOrUpdate(&entityEl)
		if err != nil {
			return
		}

		for _, name := range dataEl.Topics {
			topicEl, err := topic.FirstOrCreate(name) // TODO keep cache in memory?
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

func addTypeToTopics(data model.EntityData, t ...string) model.EntityData {
	newTopics := make([]interface{}, len(t))
	for i, v := range t {
		newTopics[i] = v
	}
	dt, ok := data["topics"]
	if !ok {
		data["topics"] = newTopics
		return data
	}

	topics, ok := dt.([]interface{})
	if !ok {
		data["topics"] = newTopics
		return data
	}

	data["topics"] = append(newTopics, topics...)
	return data
}

func useTypeAsTopic(t string) bool {
	for _, v := range autoTopicExceptions {
		if v == t {
			return false
		}
	}

	return true
}
