package timescaledb

var (
	Tab_Create_CmdHistory = `CREATE TABLE "public"."cmd_history" (
	  "id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
	  "host_id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "container_id" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "user" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "command" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "create_time" text COLLATE "pg_catalog"."default",
	  "type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'Host'::character varying
	)
	;
	ALTER TABLE "public"."cmd_history" OWNER TO "postgres";
	
	-- ----------------------------
	-- Primary Key structure for table cmd_history
	-- ----------------------------
	ALTER TABLE "public"."cmd_history" ADD CONSTRAINT "cmd_history_pkey" PRIMARY KEY ("id");`

	Tab_Create_DockerEvent = `CREATE TABLE "public"."docker_event" (
	  "id" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
	  "host_id" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "from" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "action" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "actor" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "scope" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
	  "time" int8 NOT NULL DEFAULT 0,
	  "time_nano" int8 NOT NULL DEFAULT 0
	)
	;
	ALTER TABLE "public"."docker_event" OWNER TO "postgres";
	
	-- ----------------------------
	-- Primary Key structure for table docker_event
	-- ----------------------------
	ALTER TABLE "public"."docker_event" ADD CONSTRAINT "docker_event_pkey" PRIMARY KEY ("id");`
)
