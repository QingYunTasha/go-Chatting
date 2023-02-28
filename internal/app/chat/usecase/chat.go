package usecase

/* func (u *ChatUsecase) SendUserStatus(userID uint32, status ormdomain.UserStatus) error {
	// TODO: SendStatusToPrivate()
	u.SendStatusToFriends(userID, status)
	u.SendStatusToGroup()
	return nil
}

func (u *ChatUsecase) SendGroupMessage(gorupMessage ormdomain.GroupMessage) error {
	return u.Notification.PublishMessage(context.TODO(), gorupMessage)
}

func (u *ChatUsecase) SendStatusToFriends(userID uint32, status ormdomain.UserStatus) error {
	friends, err := u.OrmRepo.User.GetFriends(userID)
	if err != nil {
		return err
	}

	for _, friend := range friends {
		conn, ok := u.ConnMgr.GetConnection(friend.ID)
		if !ok {
			continue
		}

		wsutil.WriteServerMessage(conn, ws.OpText)
	}
}
func (u *ChatUsecase) SendStatusToGroup() {}

func (u *ChatUsecase) SendUnsendMessage(conn net.Conn, messages []ormdomain.PrivateMessage) error {
	for _, msg := range messages {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err)
		}
		wsutil.WriteServerMessage(conn, ws.OpText, msgBytes)
	}
	return nil
} */

/* func (u *ChatUsecase) GetSenderGroups(SenderID uint32) ([]ormdomain.Group, error) {
	groups, err := u.OrmRepo.User.GetGroups(SenderID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
func (u *ChatUsecase) GetGroupOnlineUsers(groups []ormdomain.Group) ([]ormdomain.User, error) {
	var onlineUsers []ormdomain.User
	for _, group := range groups {
		users, err := u.OrmRepo.Group.GetUsers(group.ID)
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			if user.Status == ormdomain.Online {
				// TODO: check repeated user
				onlineUsers = append(onlineUsers, user)
			}
		}
	}

	return onlineUsers, nil
}
func (u *ChatUsecase) SendMessageToUsers(message ormdomain.GroupMessage, users []ormdomain.User) error {
	u.Notification.
}
*/
