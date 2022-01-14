package service

import (
	"78concepts.com/domicile/internal/broker"
	"78concepts.com/domicile/internal/model"
	"78concepts.com/domicile/internal/repository"
	"context"
	"encoding/json"
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

const (
	TopicGroups = broker.TopicRoot + "/bridge/groups"
)

func NewGroupsService(groupsRepository repository.IGroupsRepository) *GroupsService {
	return &GroupsService{groupsRepository: groupsRepository}
}

type GroupsService struct {
	groupsRepository repository.IGroupsRepository
}

func (s *GroupsService) ManageGroups(mqttClient *broker.MqttClient) {

	if token := mqttClient.Conn.Subscribe(TopicGroups, 0, func(client mqtt.Client, msg mqtt.Message) {
		s.HandleGroupsMessage(mqttClient.Ctx, msg);
	}); token.Wait() && token.Error() != nil {
		log.Fatal("HandleGroups: Subscribe error: %s", token.Error())
		return
	}
}

func (s *GroupsService) HandleGroupsMessage(ctx context.Context, msg mqtt.Message) {

	log.Printf("Received groups message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var objects []map[string]interface{}

	err := json.Unmarshal(msg.Payload(), &objects)

	if err != nil {
		log.Fatal(err)
	}

	groups, err:= s.GetGroups(ctx)

	for _, object := range objects {

		var found *model.Group

		for i := range groups {
			if groups[i].Id == uint64(object["id"].(float64)) {
				found = &groups[i]
				break
			}
		}

		if found != nil {

			object["active"] = true

			if !found.Active {
				object["active"] = true
				s.UpdateGroup(ctx, object)
			}

			if found.FriendlyName != object["friendly_name"] {
				s.UpdateGroup(ctx, object)
			}

		} else {
			found, _ = s.CreateGroup(ctx, object)
		}

		s.HandleGroupsMembersMessage(ctx, found, object["members"].([]interface{}))
	}

	// If a group in the database is no longer being reported, mark it as inactive
	for _, group := range groups {

		var found *map[string]interface {}

		for i := range objects {
			if uint64(objects[i]["id"].(float64)) == group.Id && group.Active {
				found = &objects[i]
				break
			}
		}

		if found == nil {

			object := map[string]interface{} {
				"id":  group.Id,
				"friendly_name":  group.FriendlyName,
				"active":  false,
			}

			s.UpdateGroup(ctx, object)
		}
	}
}

func (s *GroupsService) HandleGroupsMembersMessage(ctx context.Context, group *model.Group, objects []interface{}) {

	if group == nil {
		return
	}

	groupMembers, _ := s.GetGroupMembers(ctx, group.Id)
	for _, object := range objects {

		var found *model.GroupMember

		for i := range groupMembers {
			if groupMembers[i].IeeeAddress == object.(map[string]interface{})["ieee_address"].(string) {
				found = &groupMembers[i]
				break
			}
		}

		if found == nil {
			found, _ = s.CreateGroupMember(ctx, group.Id, object.(map[string]interface{})["ieee_address"].(string))
		}
	}

	// If a group member in the database is no longer being reported, delete it
	for _, groupMember := range groupMembers {

		var found *map[string]interface {}

		for i := range objects {
			if objects[i].(map[string]interface{})["ieee_address"].(string) == groupMember.IeeeAddress {
				foundMember := objects[i].(map[string]interface{})
				found = &foundMember
				break
			}
		}

		if found == nil {
			s.DeleteGroupMember(ctx, groupMember.GroupId, groupMember.IeeeAddress)
		}
	}
}

func (s *GroupsService) TurnGroupOn(mqttClient *broker.MqttClient, group *model.Group) {
	payload := "{\"state\": \"on\"}"

	if token := mqttClient.Conn.Publish(broker.TopicRoot + "/" + group.FriendlyName + "/set", 0, false, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (s *GroupsService) TurnGroupOff(mqttClient *broker.MqttClient, group *model.Group) {

	payload := "{\"state\": \"off\"}"

	if token := mqttClient.Conn.Publish(broker.TopicRoot + "/" + group.FriendlyName + "/set", 0, false, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (s *GroupsService) GetGroups(ctx context.Context) ([]model.Group, error) {
	return s.groupsRepository.GetGroups(ctx)
}

func (s *GroupsService) GetGroup(ctx context.Context, id uint64) (*model.Group, error) {
	return s.groupsRepository.GetGroup(ctx, id)
}

func (s *GroupsService) CreateGroup(ctx context.Context, object map[string]interface{}) (*model.Group, error) {

	if object == nil {
		return nil, errors.New("CreateGroup: Object is null")
	}

	return s.groupsRepository.CreateGroup(
		ctx,
		uint64(object["id"].(float64)),
		object["friendly_name"].(string),
	)
}

func (s *GroupsService) UpdateGroup(ctx context.Context, object map[string]interface{}) (*model.Group, error) {

	if object == nil {
		return nil, errors.New("UpdateGroup: Object is null")
	}

	return s.groupsRepository.UpdateGroup(
		ctx,
		uint64(object["id"].(float64)),
		object["friendly_name"].(string),
		object["active"].(bool),
	)
}

func (s *GroupsService) GetGroupMembers(ctx context.Context, id uint64) ([]model.GroupMember, error) {
	return s.groupsRepository.GetGroupMembers(ctx, id)
}

func (s *GroupsService) CreateGroupMember(ctx context.Context, id uint64, ieeeAddress string) (*model.GroupMember, error) {
	return s.groupsRepository.CreateGroupMember(ctx, id, ieeeAddress)
}

func (s *GroupsService) DeleteGroupMember(ctx context.Context, id uint64, ieeeAddress string) (error) {
	return s.groupsRepository.DeleteGroupMember(ctx, id, ieeeAddress)
}
