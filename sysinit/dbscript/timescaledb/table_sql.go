package timescaledb

var (
	Tab_Create_CmdHistory = `CREATE TABLE "public"."cmd_history" (
	  "id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
	  "host_id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "host_name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "container_id" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "container_name" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "user" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "command" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "create_time" int8 NOT NULL DEFAULT 0,
	  "type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'Host'::character varying
	)
	;
	ALTER TABLE "public"."cmd_history" OWNER TO "postgres";
	
	-- ----------------------------
	-- Primary Key structure for table cmd_history
	-- ----------------------------
	ALTER TABLE "public"."cmd_history" ADD CONSTRAINT "cmd_history_pkey" PRIMARY KEY ("id");`

	Tab_Create_DockerEvent = `CREATE TABLE "public"."docker_event" (
	  "id" varchar(256) COLLATE "pg_catalog"."default" NOT NULL,
	  "host_id" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "host_name" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "from" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "type" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "action" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "actor" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "status" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "scope" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "time" int8 NOT NULL DEFAULT 0,
	  "time_nano" int8 NOT NULL DEFAULT 0
	)
	;
	ALTER TABLE "public"."docker_event" OWNER TO "postgres";
	
	-- ----------------------------
	-- Primary Key structure for table docker_event
	-- ----------------------------
	ALTER TABLE "public"."docker_event" ADD CONSTRAINT "docker_event_pkey" PRIMARY KEY ("id");`

	Tab_Create_UserEvent = `CREATE TABLE "public"."user_event" (
	  "id" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
	  "user_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "account_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "raw_log" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'null'::text,
	  "model_type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "create_time" int8 NOT NULL DEFAULT 0
	)
	;
	ALTER TABLE "public"."user_event" OWNER TO "postgres";
	
	-- ----------------------------
	-- Primary Key structure for table user_event
	-- ----------------------------
	ALTER TABLE "public"."user_event" ADD CONSTRAINT "user_event_pkey" PRIMARY KEY ("id");`

	Tab_Create_TaskLog = `CREATE TABLE "public"."task_log" (
	  "id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
	  "account" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'admin'::character varying,
	  "task" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "raw_log" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "level" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'Info'::character varying,
	  "create_time" int8 NOT NULL DEFAULT 0
	)
	;
	ALTER TABLE "public"."task_log" OWNER TO "postgres";
	
	-- ----------------------------
	-- Primary Key structure for table task_log
	-- ----------------------------
	ALTER TABLE "public"."task_log" ADD CONSTRAINT "task_log_pkey" PRIMARY KEY ("id");`

	Tab_Create_WarningInfo = `CREATE TABLE "public"."warning_info" (
	  "id" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
	  "name" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "host_id" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "host_name" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "cluster" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "account" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "info" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "level" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "create_time" int8 NOT NULL DEFAULT 0,
	  "update_time" int8 NOT NULL DEFAULT 0,
	  "proposal" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "analysis" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text
	)
	;
	ALTER TABLE "public"."warning_info" OWNER TO "postgres";
	
	-- ----------------------------
	-- Primary Key structure for table warning_info
	-- ----------------------------
	ALTER TABLE "public"."warning_info" ADD CONSTRAINT "warning_info_pkey" PRIMARY KEY ("id");`

	Tab_Create_HostPackage = `CREATE TABLE "public"."host_package" (
	  "id" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
	  "name" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "host_id" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
	)
	;
	ALTER TABLE "public"."host_package" OWNER TO "postgres";
	
	-- ----------------------------
	-- Primary Key structure for table host_package
	-- ----------------------------
	ALTER TABLE "public"."host_package" ADD CONSTRAINT "host_package_pkey" PRIMARY KEY ("id");`
)
