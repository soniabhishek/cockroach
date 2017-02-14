package log_tags

import "tag"

var MESSAGE = &tag.Tag{"message"}

var USER_ID = &tag.Tag{"user_id"}
var CLIENT_ID = &tag.Tag{"client_id"}

var RESPONSE_BODY = &tag.Tag{"response_body"}
var REQUEST_BODY = &tag.Tag{"request_body"}

var HEADER = &tag.Tag{"header"}

var RECOVER = &tag.Tag{"recover"}

var FILE_NAME = &tag.Tag{"file_name"}
var FILE_PATH = &tag.Tag{"file_path"}
var ROW_NUM = &tag.Tag{"row_number"}
var COLUMN_NUM = &tag.Tag{"column_number"}

var FLU_ID = &tag.Tag{"flu_id"}
var MASTER_FLU_ID = &tag.Tag{"master_flu_id"}
var FLU_BUILD = &tag.Tag{"flu_build"}
var FLU = &tag.Tag{"flu"}
var PROJECT_ID = &tag.Tag{"project_id"}
var WORKFLOW_ID = &tag.Tag{"workflow_id"}
var PROJECT_TAG = &tag.Tag{"project_tag"}
var STEP_ID = &tag.Tag{"step_id"}

var POSTBACK_REQUEST = &tag.Tag{"postback_request"}
var POSTBACK_RESPONSE = &tag.Tag{"postback_reponse"}
var ERROR_CODE = &tag.Tag{"error_code"}
