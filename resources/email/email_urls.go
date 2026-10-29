package email

// The path template; the client fills {api_version}.
const EmailSendPath = "/api/{api_version}/admin/email/send"

// EmailPluginUUIDHeader names the plugin that owns the template identifier; a plugin_uuid body field is ignored.
const EmailPluginUUIDHeader = "plugin_uuid"
