package tgbotapi

var (
	_ Chattable = SendChecklistConfig{}
	_ Chattable = EditMessageChecklistConfig{}
	_ Chattable = SendLivePhotoConfig{}
	_ Chattable = SendMessageDraftConfig{}
	_ Chattable = SendRichMessageConfig{}
	_ Chattable = SendRichMessageDraftConfig{}
	_ Chattable = EditEphemeralMessageTextConfig{}
	_ Chattable = EditEphemeralMessageMediaConfig{}
	_ Chattable = EditEphemeralMessageCaptionConfig{}
	_ Chattable = EditEphemeralMessageReplyMarkupConfig{}
	_ Chattable = DeleteEphemeralMessageConfig{}
	_ Chattable = ApproveSuggestedPostConfig{}
	_ Chattable = DeclineSuggestedPostConfig{}
	_ Chattable = UserProfileAudiosConfig{}
	_ Chattable = UserPersonalChatMessagesConfig{}
	_ Chattable = AnswerGuestQueryConfig{}
	_ Chattable = AnswerChatJoinRequestQueryConfig{}
	_ Chattable = SendChatJoinRequestWebAppConfig{}
	_ Chattable = DeleteMessageReactionConfig{}
	_ Chattable = DeleteAllMessageReactionsConfig{}
	_ Chattable = SetPassportDataErrorsConfig{}
	_ Chattable = GetManagedBotTokenConfig{}
	_ Chattable = ReplaceManagedBotTokenConfig{}
	_ Chattable = GetManagedBotAccessSettingsConfig{}
	_ Chattable = SetManagedBotAccessSettingsConfig{}
	_ Chattable = GetMyStarBalanceConfig{}
	_ Chattable = GetBusinessAccountStarBalanceConfig{}
	_ Chattable = TransferBusinessAccountStarsConfig{}
	_ Chattable = GetUserGiftsConfig{}
	_ Chattable = GetChatGiftsConfig{}
	_ Chattable = GetBusinessAccountGiftsConfig{}
	_ Chattable = ConvertGiftToStarsConfig{}
	_ Chattable = UpgradeGiftConfig{}
	_ Chattable = TransferGiftConfig{}
	_ Chattable = GiftPremiumSubscriptionConfig{}
	_ Chattable = SetBusinessAccountNameConfig{}
	_ Chattable = SetBusinessAccountUsernameConfig{}
	_ Chattable = SetBusinessAccountBioConfig{}
	_ Chattable = SetBusinessAccountGiftSettingsConfig{}
	_ Chattable = PostStoryConfig{}
	_ Chattable = EditStoryConfig{}
	_ Chattable = RepostStoryConfig{}
	_ Chattable = DeleteStoryConfig{}
	_ Chattable = SetMyProfilePhotoConfig{}
	_ Chattable = RemoveMyProfilePhotoConfig{}
	_ Chattable = SetBusinessAccountProfilePhotoConfig{}
	_ Chattable = RemoveBusinessAccountProfilePhotoConfig{}
	_ Chattable = SavePreparedKeyboardButtonConfig{}
	_ Chattable = ChatMemberCountConfig{}
)

var (
	_ Fileable = SetMyProfilePhotoConfig{}
	_ Fileable = SendLivePhotoConfig{}
	_ Fileable = SendRichMessageConfig{}
	_ Fileable = EditMessageTextConfig{}
	_ Fileable = SendPollConfig{}
	_ Fileable = SetBusinessAccountProfilePhotoConfig{}
	_ Fileable = PostStoryConfig{}
	_ Fileable = EditStoryConfig{}
)

var (
	_ MessageId
	_ LoginUrl
	_ InlineQueryResultGif
	_ InlineQueryResultCachedGif
	_ InlineQueryResultMpeg4Gif
	_ InlineQueryResultCachedMpeg4Gif
	_ TransactionPartnerTelegramApi
	_ InputMedia = (*InputMediaLink)(nil)
	_ InputMedia = (*InputMediaVoiceNote)(nil)
)
