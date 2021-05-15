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

	Tab_Create_VirusScan = `CREATE TABLE virus_scan
		(
		id serial not null
			constraint virus_scan_pkey
				primary key,
		name text default ''::text not null,
		task_id varchar(64),
		host_id varchar(128) default ''::text not null,
		host_name varchar(128) default ''::text not null,
		image_id varchar(256) default ''::text not null,
		image_name text default ''::text not null,
		container_id varchar(256) default ''::text not null,
		container_name text default ''::text not null,
		internal_addr varchar(256) default ''::text not null,
		public_addr varchar(256) default ''::text not null,
		type varchar(32) default ''::text not null,
		raw_log text default ''::text not null,
        created_at bigint
		);
	
    alter table virus_scan owner to postgres;`

	Tab_Create_VirusRecord = `create table public.virus_record
	(
		id serial not null
			constraint virus_record_pkey
				primary key,
		virus_scan_id integer default 0 not null,
		filename text default ''::text not null,
		virus text default ''::text not null,
		database text default ''::text not null,
		type text default ''::text not null,
		size bigint default 0 not null,
		owner text default ''::text not null,
		permission bigint default 0 not null
			constraint virus_record_permission_check
				check (permission >= 0),
		modify_time bigint default 0 not null,
		create_time bigint default 0 not null
	);

	alter table public.virus_record owner to postgres;`

	Tab_Create_ImageDetail = `create table public.image_detail
	(
		id text not null
            constraint image_detail_pkey 
                 primary key,
		image_id text default ''::text not null,
		name text default ''::text not null,
		host_id text default ''::text not null,
		host_name text default ''::text not null,
		repo_tags text default ''::text not null,
		repo_digests text default ''::text not null,
		os text default ''::text not null,
		size integer default 0 not null,
		layers integer default 0 not null,
		dockerfile text default ''::text not null,
		create_time bigint default 0 not null,
		modify_time bigint default 0 not null,
		packages_json text default ''::text not null
	);

    alter table public.image_detail owner to postgres;`

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
	  "analysis" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
	  "mode" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
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
