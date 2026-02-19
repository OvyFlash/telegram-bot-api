package tgbotapi

var (
	_ Chattable = SendChecklistConfig{}
	_ Chattable = EditMessageChecklistConfig{}
	_ Chattable = SendMessageDraftConfig{}
	_ Chattable = ApproveSuggestedPostConfig{}
	_ Chattable = DeclineSuggestedPostConfig{}
	_ Chattable = UserProfileAudiosConfig{}
	_ Chattable = SetPassportDataErrorsConfig{}
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
	_ Chattable = ChatMemberCountConfig{}
)

var (
	_ Fileable = SetMyProfilePhotoConfig{}
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
)
