package usecase

type ChatUsecase struct {
}

func NewChatUsecase() *ChatUsecase {
	return &ChatUsecase{}
}

func (u *ChatUsecase) SendMessage(senderId, receiverId, channelType, message string) error {
	return nil
}

func (u *ChatUsecase) GetMessage(senderId, receiverId, channelType, perMessageId string) error {
	return nil
}

func (u *ChatUsecase) Join(userId, GroupId string) error {
	return nil
}

func (u *ChatUsecase) LeaveGroup(userId, GroupId string) error {
	return nil
}

func (u *ChatUsecase) GetGroupsByUser(userId string) error {
	return nil
}
