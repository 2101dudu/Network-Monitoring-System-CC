package utils

type request_type byte

const (
     AGENT_REGISTRATION request_type = iota  // iota = 0
     TASK_REQUEST  // iota = 1
     METRICS_GATHERING  // iota = 2
     ERROR  // iota = 3
)


type ack struct {
    acknowledged bool
    sender_id byte  // [0, 255]
    sender_is_server bool
    request_id byte  // [0, 255]
    request_type request_type  //(AGENT_REGISTRATION, TASK_DELEGATION, DATA_COLLECTION)
}

type ack_builder struct {
    ack ack
}

func New_ack_builder() *ack_builder {
    return &ack_builder{
        ack: ack{
        acknowledged: false,
        sender_id: 0,
        sender_is_server: false,
        request_id: 0,
        request_type: ERROR},
    }
}

func (a *ack_builder) Has_ackowledged() *ack_builder {
    a.ack.acknowledged = true
    return a
}

func (a *ack_builder) Set_sender_id (id byte) *ack_builder {
    a.ack.sender_id = id
    return a
}

func (a *ack_builder) Is_server() *ack_builder {
    a.ack.sender_is_server = true
    return a
}

func (a *ack_builder) Set_request_id(id byte) *ack_builder {
    a.ack.request_id = id
    return a
}

func (a *ack_builder) Set_request_type(request request_type) *ack_builder {
    a.ack.request_type = request
    return a
}

func (a *ack_builder) Build() ack {
    return a.ack
}
