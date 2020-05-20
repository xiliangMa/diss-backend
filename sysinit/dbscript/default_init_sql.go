package dbscript

var (
	DefaultDockerBenchSql     = `INSERT INTO "public"."system_template" VALUES ('9b4320a4-57f5-4559-a91b-d84b0d0a0e99', 'admin', 'CIS标准-Docker Benchmark', '', 'DockerBenchMark', 'v.1.3.5', 'docker run --rm --net host --pid host --cap-add audit_control     -e DOCKER_CONTENT_TRUST=$DOCKER_CONTENT_TRUST      -v /var/lib:/var/lib:ro     -v /var/run/docker.sock:/var/run/docker.sock:ro       -v /usr/lib/systemd:/usr/lib/systemd:ro     -v /etc:/etc:ro    --label docker_bench_security  docker-bench-security  -o json', 'Enable', 't');`
	DefaultK8sBenchSql        = `INSERT INTO "public"."system_template" VALUES ('9b4320a4-57f5-4559-a91b-d84b0d0a0e98', 'admin', 'CIS标准-Kubernetes benchmark','', 'KubernetesBenchMark', 'cis-1.3', 'docker run --rm --pid=host -v /etc:/etc:ro -v /var:/var:ro -t aquasec/kube-bench:latest  --benchmark cis-1.3 --json', 'Enable', 't');`
	DefaultDockerVirusScanSql = `INSERT INTO "public"."system_template" VALUES ('9b4320a4-57f5-4559-a91b-d84b0d0a0e90', 'admin', '容器病毒扫描', '', 'DockerVirusScan','1.0', '', 'Enable', 't')`
    DefaultLogConfigSql = `INSERT INTO public.log_config (id, enabled, server_url, server_port, exported_types, config_name) VALUES ('cf110d02-02d9-4ab4-91cc-f4003840d649', true, '111.229.167.6', '514', 'ImageSecurityLog,BenchmarkScanLog,IntrudeDetectLog', 'SysLogExport');`
)
