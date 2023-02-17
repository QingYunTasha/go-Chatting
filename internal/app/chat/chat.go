package usecase

type ChatUsecase struct {
}

func NewChatUsecase() *ChatUsecase {
	return &ChatUsecase{}
}

func (u *ChatUsecase) SendMessage(senderID, receiverID, channelType, message string) error {
	return nil
}

func (u *ChatUsecase) GetMessage(user1ID, user2ID, channelType, preMessageID string) error {
	return nil
}

func GetStatus() {}

func SendStatus() {}
