package message

//MsgType define message type
type MsgType int8

const (
	//MSG_NAME register a name to join the hub
	MSG_NAME = 1
	//MSG_TEXT indicate this is real message
	MSG_TEXT = 2
)

type MessageRegister struct {
	Type MsgType
	Name string
}

func (m MessageRegister) Seriallize() []byte {
	return nil
}

type Message struct {
	Type MsgType
	Text string
}

func (m Message) Seriallize() []byte {
	return nil
}

type MessageBroadcast struct {
	Type MsgType
	From string
	Text string
}

func (m MessageBroadcast) Seriallize() []byte {
	var fromBytes = []byte(m.From)
	var data = append([]byte{byte(m.Type)}, []byte{byte(len(fromBytes))}...)
	data = append(data, fromBytes...)
	data = append(data, []byte(m.Text)...)
	return data
}

/*


 */
