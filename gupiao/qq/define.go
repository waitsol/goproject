package qq

const (
	GUILDS         = 1 << iota // 1 << 0
	GUILD_CREATE               // 当机器人加入新guild时
	GUILD_UPDATE               // 当guild资料发生变更时
	GUILD_DELETE               // 当机器人退出guild时
	CHANNEL_CREATE             // 当channel被创建时
	CHANNEL_UPDATE             // 当channel被更新时
	CHANNEL_DELETE             // 当channel被删除时

	GUILD_MEMBERS       = 1 << 1 // 1 << 1
	GUILD_MEMBER_ADD             // 当成员加入时
	GUILD_MEMBER_UPDATE          // 当成员资料变更时
	GUILD_MEMBER_REMOVE          // 当成员被移除时

	GUILD_MESSAGES = 1 << 9 // 1 << 9
	MESSAGE_CREATE          // 发送消息事件
	MESSAGE_DELETE          // 删除（撤回）消息事件

	GUILD_MESSAGE_REACTIONS = 1 << 10 // 1 << 10
	MESSAGE_REACTION_ADD              // 为消息添加表情表态
	MESSAGE_REACTION_REMOVE           // 为消息删除表情表态

	DIRECT_MESSAGE        = 1 << 12 // 1 << 12
	DIRECT_MESSAGE_CREATE           // 收到用户发给机器人的私信消息
	DIRECT_MESSAGE_DELETE           // 删除（撤回）消息事件

	GROUP_AND_C2C_EVENT     = 1 << 25 // 1 << 25
	C2C_MESSAGE_CREATE                // 用户单聊发消息给机器人时候
	FRIEND_ADD                        // 用户添加使用机器人
	FRIEND_DEL                        // 用户删除机器人
	C2C_MSG_REJECT                    // 用户手动关闭"主动消息"推送
	C2C_MSG_RECEIVE                   // 用户手动开启"主动消息"推送
	GROUP_AT_MESSAGE_CREATE           // 用户在群里@机器人时收到的消息
	GROUP_ADD_ROBOT                   // 机器人被添加到群聊
	GROUP_DEL_ROBOT                   // 机器人被移出群聊
	GROUP_MSG_REJECT                  // 群管理员关闭通知
	GROUP_MSG_RECEIVE                 // 群管理员开启通知

	INTERACTION        = 1 << 26 // 1 << 26
	INTERACTION_CREATE           // 互动事件创建时

	MESSAGE_AUDIT        = 1 << 27 // 1 << 27
	MESSAGE_AUDIT_PASS             // 消息审核通过
	MESSAGE_AUDIT_REJECT           // 消息审核不通过

	FORUMS_EVENT               = 1 << 28 // 1 << 28
	FORUM_THREAD_CREATE                  // 当用户创建主题时
	FORUM_THREAD_UPDATE                  // 当用户更新主题时
	FORUM_THREAD_DELETE                  // 当用户删除主题时
	FORUM_POST_CREATE                    // 当用户创建帖子时
	FORUM_POST_DELETE                    // 当用户删除帖子时
	FORUM_REPLY_CREATE                   // 当用户回复评论时
	FORUM_REPLY_DELETE                   // 当用户回复评论时
	FORUM_PUBLISH_AUDIT_RESULT           // 当用户发表审核通过时

	AUDIO_ACTION  = 1 << 29 // 1 << 29
	AUDIO_START             // 音频开始播放时
	AUDIO_FINISH            // 音频播放结束时
	AUDIO_ON_MIC            // 上麦时
	AUDIO_OFF_MIC           // 下麦时

	PUBLIC_GUILD_MESSAGES = 1 << 30 // 1 << 30
	AT_MESSAGE_CREATE               // 收到@机器人的消息时
	PUBLIC_MESSAGE_DELETE           // 当频道的消息被删除时
)
