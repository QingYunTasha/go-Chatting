package usecase

type ChatUsecase struct {
}

func NewChatUsecase() *ChatUsecase {
	return &ChatUsecase{}
}

func (u *ChatUsecase) SendMessage(senderID, receiverID, channelType, message string) error {
	return nil
}

func (u *ChatUsecase) GetMessage(senderID, receiverID, channelType, preMessageID string) error {
	return nil
}

func GetStatus(senderID, receiverID, channelType, status string) error {
	return nil
}

func SendStatus(senderID, receiverID, channelType, status string) error {
	return nil
}
