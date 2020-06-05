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
)
