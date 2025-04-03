package rscp

import (
	"fmt"
)

// validateResponse checks the integrity of the response
// must contain a valid tag and data type and the data type must match the value
func validateResponse(message Message) error {
	if !message.Tag.isResponse() {
		return fmt.Errorf("%s: %w", message.Tag, ErrNotAResponseTag)
	}
	return message.validate()
}

// ValidateResponses checks the integrity of the responses
// each response must contain a valid tag and data type and the data type must match the value
func ValidateResponses(messages []Message) error {
	for i, m := range messages {
		if err := validateResponse(m); err != nil {
			return fmt.Errorf("message at index %d: %w", i, err)
		}
	}
	return nil
}

// CreateResponse creates a new message (infer the data type from the tag)
// if the tag's data type is a Container, following tag's will be nested as "sub" responses within the container.
// Every tag that has a data type other than DATATYPE_None requires a following value.
// Examples:
//
//	CreateResponse(INFO_REQ_UTC_TIME)
//	CreateResponse(EMS_REQ_SET_ERROR_BUZZER_ENABLED, true)
//	CreateResponse(BAT_REQ_DATA, BAT_INDEX, uint16(0), BAT_REQ_DEVICE_STATE, BAT_REQ_RSOC, BAT_REQ_STATUS_CODE)
func CreateResponse(values ...interface{}) (msg *Message, err error) {
	if msg, err = readResponseSlice(values); err != nil {
		return nil, err
	}
	return msg, nil
}

// CreateResponses creates multiple new responses (infer the data type from the tag)
// if the tag's data type is a Container, provided values will be converted to "sub" responses, separated by the provided tag's
// Examples:
//
//	CreateResponses([]interface{}{INFO_REQ_UTC_TIME})
//	CreateResponses([]interface{}{EMS_REQ_SET_ERROR_BUZZER_ENABLED, true})
//	CreateResponses([]interface{}{BAT_REQ_DATA, BAT_INDEX, uint16(0), BAT_REQ_DEVICE_STATE, BAT_REQ_RSOC, BAT_REQ_STATUS_CODE})
func CreateResponses(values ...[]interface{}) ([]Message, error) {
	if len(values) == 0 {
		return nil, ErrNoArguments
	}
	msgs := make([]Message, 0)
	for _, subValues := range values {
		var (
			msg *Message
			err error
		)
		if msg, err = CreateResponse(subValues...); err != nil {
			return nil, err
		}
		msgs = append(msgs, *msg)
	}
	return msgs, nil
}
