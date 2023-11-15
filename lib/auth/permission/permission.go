package permission

import "github.com/bwmarrin/discordgo"

const (
	PermissionSendVoiceMessages    = 0x0000400000000000
	PermissionVcUseSoundboard      = 0x0000040000000000
	PermissionVcUserExternalSounds = 0x0000200000000000
	PermissionCreateExpressions    = 0x80000000000
	PermissionCreateEvents         = 0x100000000000
)

const (
	ChannelTypeText     = "text"
	ChannelTypeCategory = "category"
	ChannelTypeAnnounce = "announce"
	ChannelTypeForum    = "forum"
	ChannelTypeVC       = "vc"
	ChannelTypeStage    = "stage"
)

// Typeを変換します
func ConvertChannelType(chType discordgo.ChannelType) string {
	switch chType {
	case discordgo.ChannelTypeGuildText:
		return ChannelTypeText
	case discordgo.ChannelTypeGuildCategory:
		return ChannelTypeCategory
	case discordgo.ChannelTypeGuildNews:
		return ChannelTypeAnnounce
	case discordgo.ChannelTypeGuildForum:
		return ChannelTypeForum
	case discordgo.ChannelTypeGuildVoice:
		return ChannelTypeVC
	case discordgo.ChannelTypeGuildStageVoice:
		return ChannelTypeStage
	default:
		return ""
	}
}

type ViewChannels bool
type ManageChannels bool
type ManageRoles bool
type CreateExpressions bool
type ManageExpressions bool // 絵文字の管理
type ViewAuditLog bool
type ViewServerInsights bool
type ManageWebhooks bool
type ManageServer bool
type CreateInvite bool
type ChangeNickname bool
type ManageNickname bool
type KickMembers bool
type BanMembers bool
type TimeoutMembers bool
type SendMessages bool
type SendMessagesInThreads bool
type CreatePublicThreads bool
type CreatePrivateThreads bool
type EmbedLinks bool
type AttachFiles bool
type AddReactions bool
type UseExternalEmoji bool
type UseExternalStickers bool
type MentionEveryone bool
type ManageMessages bool
type ManageThreads bool
type ReadMessageHistory bool
type SendTextToSpeechMessage bool
type UseApplicationCommands bool
type SendVoiceMessages bool
type VcConnect bool
type VcSpeak bool
type VcVideo bool
type VcUseActivities bool
type VcUseSoundboard bool
type VcUseExternalSounds bool
type VcUseVoiceActivity bool
type VcPrioritySpeaker bool
type VcMuteMembers bool
type VcDeafenMembers bool
type VcMoveMembers bool
type StageRequestToSpeak bool
type CreateEvents bool
type ManageEvents bool
type Administrator bool

type Permission interface {
	Permission()
}

type RolePermission struct {
	ViewChannels            `json:"view_channels"`
	ManageChannels          `json:"manage_channels"`
	ManageRoles             `json:"manage_roles"`
	CreateExpressions       `json:"create_expressions"`
	ManageExpressions       `json:"manage_expressions"`
	ViewAuditLog            `json:"view_audit_log"`
	ViewServerInsights      `json:"view_server_insights"`
	ManageWebhooks          `json:"manage_webhooks"`
	ManageServer            `json:"manage_server"`
	CreateInvite            `json:"create_invite"`
	ChangeNickname          `json:"change_nickname"`
	ManageNickname          `json:"manage_nickname"`
	KickMembers             `json:"kick_members"`
	BanMembers              `json:"ban_members"`
	TimeoutMembers          `json:"timeout_members"`
	SendMessages            `json:"send_messages"`
	SendMessagesInThreads   `json:"send_messages_in_threads"`
	CreatePublicThreads     `json:"create_public_threads"`
	CreatePrivateThreads    `json:"create_private_threads"`
	EmbedLinks              `json:"embed_links"`
	AttachFiles             `json:"attach_files"`
	AddReactions            `json:"add_reactions"`
	UseExternalEmoji        `json:"use_external_emoji"`
	UseExternalStickers     `json:"use_external_stickers"`
	MentionEveryone         `json:"mention_everyone"`
	ManageMessages          `json:"manage_messages"`
	ManageThreads           `json:"manage_threads"`
	ReadMessageHistory      `json:"read_message_history"`
	SendTextToSpeechMessage `json:"send_text_to_speech_message"`
	UseApplicationCommands  `json:"use_application_commands"`
	SendVoiceMessages       `json:"send_voice_messages"`
	VcConnect               `json:"vc_connect"`
	VcSpeak                 `json:"vc_speak"`
	VcVideo                 `json:"vc_video"`
	VcUseActivities         `json:"vc_use_activities"`
	VcUseSoundboard         `json:"vc_use_soundboard"`
	VcUseExternalSounds     `json:"vc_use_external_sounds"`
	VcUseVoiceActivity      `json:"vc_use_voice_activity"`
	VcPrioritySpeaker       `json:"vc_priority_speaker"`
	VcMuteMembers           `json:"vc_mute_members"`
	VcDeafenMembers         `json:"vc_deafen_members"`
	VcMoveMembers           `json:"vc_move_members"`
	StageRequestToSpeak     `json:"stage_request_to_speak"`
	CreateEvents            `json:"create_events"`
	ManageEvents            `json:"manage_events"`
	Administrator           `json:"administrator"`
}

func (p RolePermission) Permission() {}

type TextChannelPermission struct {
	Administrator           `json:"administrator"`
	ViewChannels            `json:"view_channels"`
	ManageChannels          `json:"manage_channels"`
	ManageRoles             `json:"manage_roles"`
	ManageWebhooks          `json:"manage_webhooks"`
	CreateInvite            `json:"create_invite"`
	SendMessages            `json:"send_messages"`
	SendMessagesInThreads   `json:"send_messages_in_threads"`
	CreatePublicThreads     `json:"create_public_threads"`
	CreatePrivateThreads    `json:"create_private_threads"`
	EmbedLinks              `json:"embed_links"`
	AttachFiles             `json:"attach_files"`
	AddReactions            `json:"add_reactions"`
	UseExternalEmoji        `json:"use_external_emoji"`
	UseExternalStickers     `json:"use_external_stickers"`
	MentionEveryone         `json:"mention_everyone"`
	ManageMessages          `json:"manage_messages"`
	ManageThreads           `json:"manage_threads"`
	ReadMessageHistory      `json:"read_message_history"`
	SendTextToSpeechMessage `json:"send_text_to_speech_message"`
	UseApplicationCommands  `json:"use_application_commands"`
	SendVoiceMessages       `json:"send_voice_messages"`
}

func (p TextChannelPermission) Permission() {}

type CategoryPermission struct {
	Administrator           `json:"administrator"`
	ViewChannels            `json:"view_channels"`
	ManageChannels          `json:"manage_channels"`
	ManageRoles             `json:"manage_roles"`
	ManageWebhooks          `json:"manage_webhooks"`
	CreateInvite            `json:"create_invite"`
	SendMessages            `json:"send_messages"`
	SendMessagesInThreads   `json:"send_messages_in_threads"`
	CreatePublicThreads     `json:"create_public_threads"`
	CreatePrivateThreads    `json:"create_private_threads"`
	EmbedLinks              `json:"embed_links"`
	AttachFiles             `json:"attach_files"`
	AddReactions            `json:"add_reactions"`
	UseExternalEmoji        `json:"use_external_emoji"`
	UseExternalStickers     `json:"use_external_stickers"`
	MentionEveryone         `json:"mention_everyone"`
	ManageMessages          `json:"manage_messages"`
	ManageThreads           `json:"manage_threads"`
	ReadMessageHistory      `json:"read_message_history"`
	SendTextToSpeechMessage `json:"send_text_to_speech_message"`
	UseApplicationCommands  `json:"use_application_commands"`
	SendVoiceMessages       `json:"send_voice_messages"`
	VcConnect               `json:"vc_connect"`
	VcSpeak                 `json:"vc_speak"`
	VcVideo                 `json:"vc_video"`
	VcUseActivities         `json:"vc_use_activities"`
	VcUseSoundboard         `json:"vc_use_soundboard"`
	VcUseExternalSounds     `json:"vc_use_external_sounds"`
	VcUseVoiceActivity      `json:"vc_use_voice_activity"`
	VcPrioritySpeaker       `json:"vc_priority_speaker"`
	VcMuteMembers           `json:"vc_mute_members"`
	VcDeafenMembers         `json:"vc_deafen_members"`
	VcMoveMembers           `json:"vc_move_members"`
	StageRequestToSpeak     `json:"stage_request_to_speak"`
	CreateEvents            `json:"create_events"`
	ManageEvents            `json:"manage_events"`
}

func (p CategoryPermission) Permission() {}

type AnnounceChannelPermission struct {
	Administrator           `json:"administrator"`
	ViewChannels            `json:"view_channels"`
	ManageChannels          `json:"manage_channels"`
	ManageRoles             `json:"manage_roles"`
	ManageWebhooks          `json:"manage_webhooks"`
	CreateInvite            `json:"create_invite"`
	SendMessages            `json:"send_messages"`
	SendMessagesInThreads   `json:"send_messages_in_threads"`
	CreatePublicThreads     `json:"create_public_threads"`
	EmbedLinks              `json:"embed_links"`
	AttachFiles             `json:"attach_files"`
	AddReactions            `json:"add_reactions"`
	UseExternalEmoji        `json:"use_external_emoji"`
	UseExternalStickers     `json:"use_external_stickers"`
	MentionEveryone         `json:"mention_everyone"`
	ManageMessages          `json:"manage_messages"`
	ManageThreads           `json:"manage_threads"`
	ReadMessageHistory      `json:"read_message_history"`
	SendTextToSpeechMessage `json:"send_text_to_speech_message"`
	UseApplicationCommands  `json:"use_application_commands"`
	SendVoiceMessages       `json:"send_voice_messages"`
}

func (p AnnounceChannelPermission) Permission() {}

type ForumPermission struct {
	Administrator           `json:"administrator"`
	ViewChannels            `json:"view_channels"`
	ManageChannels          `json:"manage_channels"`
	ManageRoles             `json:"manage_roles"`
	ManageWebhooks          `json:"manage_webhooks"`
	CreateInvite            `json:"create_invite"`
	SendMessages            `json:"send_messages"`
	SendMessagesInThreads   `json:"send_messages_in_threads"`
	EmbedLinks              `json:"embed_links"`
	AttachFiles             `json:"attach_files"`
	AddReactions            `json:"add_reactions"`
	UseExternalEmoji        `json:"use_external_emoji"`
	UseExternalStickers     `json:"use_external_stickers"`
	MentionEveryone         `json:"mention_everyone"`
	ManageMessages          `json:"manage_messages"`
	ManageThreads           `json:"manage_threads"`
	ReadMessageHistory      `json:"read_message_history"`
	SendTextToSpeechMessage `json:"send_text_to_speech_message"`
	UseApplicationCommands  `json:"use_application_commands"`
	SendVoiceMessages       `json:"send_voice_messages"`
}

func (p ForumPermission) Permission() {}

type VCPermission struct {
	Administrator           `json:"administrator"`
	ViewChannels            `json:"view_channels"`
	ManageChannels          `json:"manage_channels"`
	ManageRoles             `json:"manage_roles"`
	ManageWebhooks          `json:"manage_webhooks"`
	CreateInvite            `json:"create_invite"`
	SendMessages            `json:"send_messages"`
	EmbedLinks              `json:"embed_links"`
	AttachFiles             `json:"attach_files"`
	AddReactions            `json:"add_reactions"`
	UseExternalEmoji        `json:"use_external_emoji"`
	UseExternalStickers     `json:"use_external_stickers"`
	MentionEveryone         `json:"mention_everyone"`
	ManageMessages          `json:"manage_messages"`
	ReadMessageHistory      `json:"read_message_history"`
	SendTextToSpeechMessage `json:"send_text_to_speech_message"`
	UseApplicationCommands  `json:"use_application_commands"`
	VcConnect               `json:"vc_connect"`
	VcSpeak                 `json:"vc_speak"`
	VcVideo                 `json:"vc_video"`
	VcUseActivities         `json:"vc_use_activities"`
	VcUseSoundboard         `json:"vc_use_soundboard"`
	VcUseExternalSounds     `json:"vc_use_external_sounds"`
	VcUseVoiceActivity      `json:"vc_use_voice_activity"`
	VcPrioritySpeaker       `json:"vc_priority_speaker"`
	VcMuteMembers           `json:"vc_mute_members"`
	VcDeafenMembers         `json:"vc_deafen_members"`
	VcMoveMembers           `json:"vc_move_members"`
	ManageEvents            `json:"manage_events"`
}

func (p VCPermission) Permission() {}

type StagePermission struct {
	Administrator           `json:"administrator"`
	ViewChannels            `json:"view_channels"`
	ManageChannels          `json:"manage_channels"`
	ManageRoles             `json:"manage_roles"`
	CreateInvite            `json:"create_invite"`
	SendMessages            `json:"send_messages"`
	EmbedLinks              `json:"embed_links"`
	AttachFiles             `json:"attach_files"`
	AddReactions            `json:"add_reactions"`
	UseExternalEmoji        `json:"use_external_emoji"`
	UseExternalStickers     `json:"use_external_stickers"`
	MentionEveryone         `json:"mention_everyone"`
	ManageMessages          `json:"manage_messages"`
	ReadMessageHistory      `json:"read_message_history"`
	SendTextToSpeechMessage `json:"send_text_to_speech_message"`
	UseApplicationCommands  `json:"use_application_commands"`
	VcConnect               `json:"vc_connect"`
	VcVideo                 `json:"vc_video"`
	VcMuteMembers           `json:"vc_mute_members"`
	VcMoveMembers           `json:"vc_move_members"`
	StageRequestToSpeak     `json:"stage_request_to_speak"`
	ManageEvents            `json:"manage_events"`
}

func (p StagePermission) Permission() {}
