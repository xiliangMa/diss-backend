package models

import (
	"encoding/json"
	"testing"
)

func TestDecodeImageVulnerabilities(t *testing.T) {
	imageVulnerabilitiesStr := `[
  {
    "Target": "python:3.4-alpine (alpine 3.9.2)",
    "Type": "alpine",
    "Vulnerabilities": [
      {
        "VulnerabilityID": "CVE-2018-20843",
        "PkgName": "expat",
        "InstalledVersion": "2.2.6-r0",
        "FixedVersion": "2.2.7-r0",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2018-20843",
        "Title": "expat: large number of colons in input makes parser consume high amount of resources, leading to DoS",
        "Description": "In libexpat in Expat before 2.2.7, XML input including XML names that contain a large number of colons could make the XML parser consume a high amount of RAM and CPU resources while processing (enough to be usable for denial-of-service attacks).",
        "Severity": "HIGH",
        "CweIDs": [
          "CWE-611"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:N/I:N/A:C",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 7.8,
            "V3Score": 7.5
          },
          "redhat": {
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V3Score": 7.5
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-07/msg00039.html",
          "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=5226",
          "https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=931031",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2018-20843",
          "https://github.com/libexpat/libexpat/blob/R_2_2_7/expat/Changes",
          "https://github.com/libexpat/libexpat/issues/186",
          "https://github.com/libexpat/libexpat/pull/262",
          "https://github.com/libexpat/libexpat/pull/262/commits/11f8838bf99ea0a6f0b76f9760c43704d00c4ff6",
          "https://linux.oracle.com/cve/CVE-2018-20843.html",
          "https://linux.oracle.com/errata/ELSA-2020-4484.html",
          "https://lists.debian.org/debian-lts-announce/2019/06/msg00028.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/CEJJSQSG3KSUQY4FPVHZ7ZTT7FORMFVD/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/IDAUGEB3TUP6NEKJDBUBZX7N5OAUOOOK/",
          "https://seclists.org/bugtraq/2019/Jun/39",
          "https://security.gentoo.org/glsa/201911-08",
          "https://security.netapp.com/advisory/ntap-20190703-0001/",
          "https://support.f5.com/csp/article/K51011533",
          "https://usn.ubuntu.com/4040-1/",
          "https://usn.ubuntu.com/4040-2/",
          "https://usn.ubuntu.com/usn/usn-4040-1",
          "https://usn.ubuntu.com/usn/usn-4040-2",
          "https://www.debian.org/security/2019/dsa-4472",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html"
        ],
        "PublishedDate": "2019-06-24T17:15:00Z",
        "LastModifiedDate": "2021-01-25T18:12:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-15903",
        "PkgName": "expat",
        "InstalledVersion": "2.2.6-r0",
        "FixedVersion": "2.2.7-r1",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-15903",
        "Title": "expat: heap-based buffer over-read via crafted XML input",
        "Description": "In libexpat before 2.2.8, crafted XML input could fool the parser into changing from DTD parsing to document parsing too early; a consecutive call to XML_GetCurrentLineNumber (or XML_GetCurrentColumnNumber) then resulted in a heap-based buffer over-read.",
        "Severity": "HIGH",
        "CweIDs": [
          "CWE-776",
          "CWE-125"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 5,
            "V3Score": 7.5
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V3Score": 7.5
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00080.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00081.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00000.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00002.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00003.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00013.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00016.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00017.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00018.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00019.html",
          "http://lists.opensuse.org/opensuse-security-announce/2020-01/msg00008.html",
          "http://lists.opensuse.org/opensuse-security-announce/2020-01/msg00040.html",
          "http://packetstormsecurity.com/files/154503/Slackware-Security-Advisory-expat-Updates.html",
          "http://packetstormsecurity.com/files/154927/Slackware-Security-Advisory-python-Updates.html",
          "http://packetstormsecurity.com/files/154947/Slackware-Security-Advisory-mozilla-firefox-Updates.html",
          "http://seclists.org/fulldisclosure/2019/Dec/23",
          "http://seclists.org/fulldisclosure/2019/Dec/26",
          "http://seclists.org/fulldisclosure/2019/Dec/27",
          "http://seclists.org/fulldisclosure/2019/Dec/30",
          "https://access.redhat.com/errata/RHSA-2019:3210",
          "https://access.redhat.com/errata/RHSA-2019:3237",
          "https://access.redhat.com/errata/RHSA-2019:3756",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-15903",
          "https://github.com/libexpat/libexpat/commit/c20b758c332d9a13afbbb276d30db1d183a85d43",
          "https://github.com/libexpat/libexpat/issues/317",
          "https://github.com/libexpat/libexpat/issues/342",
          "https://github.com/libexpat/libexpat/pull/318",
          "https://linux.oracle.com/cve/CVE-2019-15903.html",
          "https://linux.oracle.com/errata/ELSA-2020-4484.html",
          "https://lists.debian.org/debian-lts-announce/2019/11/msg00006.html",
          "https://lists.debian.org/debian-lts-announce/2019/11/msg00017.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/A4TZKPJFTURRLXIGLB34WVKQ5HGY6JJA/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/BDUTI5TVQWIGGQXPEVI4T2ENHFSBMIBP/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/S26LGXXQ7YF2BP3RGOWELBFKM6BHF6UG/",
          "https://seclists.org/bugtraq/2019/Dec/17",
          "https://seclists.org/bugtraq/2019/Dec/21",
          "https://seclists.org/bugtraq/2019/Dec/23",
          "https://seclists.org/bugtraq/2019/Nov/1",
          "https://seclists.org/bugtraq/2019/Nov/24",
          "https://seclists.org/bugtraq/2019/Oct/29",
          "https://seclists.org/bugtraq/2019/Sep/30",
          "https://seclists.org/bugtraq/2019/Sep/37",
          "https://security.gentoo.org/glsa/201911-08",
          "https://security.netapp.com/advisory/ntap-20190926-0004/",
          "https://support.apple.com/kb/HT210785",
          "https://support.apple.com/kb/HT210788",
          "https://support.apple.com/kb/HT210789",
          "https://support.apple.com/kb/HT210790",
          "https://support.apple.com/kb/HT210793",
          "https://support.apple.com/kb/HT210794",
          "https://support.apple.com/kb/HT210795",
          "https://usn.ubuntu.com/4132-1/",
          "https://usn.ubuntu.com/4132-2/",
          "https://usn.ubuntu.com/4165-1/",
          "https://usn.ubuntu.com/4202-1/",
          "https://usn.ubuntu.com/4335-1/",
          "https://usn.ubuntu.com/usn/usn-4132-1",
          "https://usn.ubuntu.com/usn/usn-4132-2",
          "https://usn.ubuntu.com/usn/usn-4165-1",
          "https://usn.ubuntu.com/usn/usn-4202-1",
          "https://usn.ubuntu.com/usn/usn-4335-1",
          "https://www.debian.org/security/2019/dsa-4530",
          "https://www.debian.org/security/2019/dsa-4549",
          "https://www.debian.org/security/2019/dsa-4571",
          "https://www.mozilla.org/en-US/security/advisories/mfsa2019-34/#CVE-2019-15903",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html"
        ],
        "PublishedDate": "2019-09-04T06:15:00Z",
        "LastModifiedDate": "2020-10-20T22:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-12900",
        "PkgName": "libbz2",
        "InstalledVersion": "1.0.6-r6",
        "FixedVersion": "1.0.6-r7",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-12900",
        "Title": "bzip2: out-of-bounds write in function BZ2_decompress",
        "Description": "BZ2_decompress in decompress.c in bzip2 through 1.0.6 has an out-of-bounds write when there are many selectors.",
        "Severity": "CRITICAL",
        "CweIDs": [
          "CWE-787"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:P/I:P/A:P",
            "V3Vector": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
            "V2Score": 7.5,
            "V3Score": 9.8
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:L/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:L",
            "V3Score": 4
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-07/msg00040.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-08/msg00050.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-11/msg00078.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-12/msg00000.html",
          "http://packetstormsecurity.com/files/153644/Slackware-Security-Advisory-bzip2-Updates.html",
          "http://packetstormsecurity.com/files/153957/FreeBSD-Security-Advisory-FreeBSD-SA-19-18.bzip2.html",
          "https://bugs.launchpad.net/ubuntu/+source/bzip2/+bug/1834494",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-12900",
          "https://gitlab.com/federicomenaquintero/bzip2/commit/74de1e2e6ffc9d51ef9824db71a8ffee5962cdbc",
          "https://lists.apache.org/thread.html/ra0adb9653c7de9539b93cc8434143b655f753b9f60580ff260becb2b@%3Cusers.kafka.apache.org%3E",
          "https://lists.debian.org/debian-lts-announce/2019/06/msg00021.html",
          "https://lists.debian.org/debian-lts-announce/2019/07/msg00014.html",
          "https://lists.debian.org/debian-lts-announce/2019/10/msg00012.html",
          "https://lists.debian.org/debian-lts-announce/2019/10/msg00018.html",
          "https://seclists.org/bugtraq/2019/Aug/4",
          "https://seclists.org/bugtraq/2019/Jul/22",
          "https://security.FreeBSD.org/advisories/FreeBSD-SA-19:18.bzip2.asc",
          "https://support.f5.com/csp/article/K68713584?utm_source=f5support\u0026amp;utm_medium=RSS",
          "https://usn.ubuntu.com/4038-1/",
          "https://usn.ubuntu.com/4038-2/",
          "https://usn.ubuntu.com/4146-1/",
          "https://usn.ubuntu.com/4146-2/",
          "https://usn.ubuntu.com/usn/usn-4038-1",
          "https://usn.ubuntu.com/usn/usn-4038-2",
          "https://usn.ubuntu.com/usn/usn-4038-3",
          "https://usn.ubuntu.com/usn/usn-4038-4",
          "https://usn.ubuntu.com/usn/usn-4146-1",
          "https://usn.ubuntu.com/usn/usn-4146-2",
          "https://www.oracle.com/security-alerts/cpuoct2020.html"
        ],
        "PublishedDate": "2019-06-19T23:15:00Z",
        "LastModifiedDate": "2020-10-20T22:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1543",
        "PkgName": "libcrypto1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1b-r1",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1543",
        "Title": "openssl: ChaCha20-Poly1305 with long nonces",
        "Description": "ChaCha20-Poly1305 is an AEAD cipher, and requires a unique nonce input for every encryption operation. RFC 7539 specifies that the nonce value (IV) should be 96 bits (12 bytes). OpenSSL allows a variable nonce length and front pads the nonce with 0 bytes if it is less than 12 bytes. However it also incorrectly allows a nonce to be set of up to 16 bytes. In this case only the last 12 bytes are significant and any additional leading bytes are ignored. It is a requirement of using this cipher that nonce values are unique. Messages encrypted using a reused nonce value are susceptible to serious confidentiality and integrity attacks. If an application changes the default nonce length to be longer than 12 bytes and then makes a change to the leading bytes of the nonce expecting the new value to be a new unique nonce then such an application could inadvertently encrypt messages with a reused nonce. Additionally the ignored bytes in a long nonce are not covered by the integrity guarantee of this cipher. Any application that relies on the integrity of these ignored leading bytes of a long nonce may be further affected. Any OpenSSL internal use of this cipher, including in SSL/TLS, is safe because no such use sets such a long nonce value. However user applications that use this cipher directly and set a non-default nonce length to be longer than 12 bytes may be vulnerable. OpenSSL versions 1.1.1 and 1.1.0 are affected by this issue. Due to the limited scope of affected deployments this has been assessed as low severity and therefore we are not creating new releases at this time. Fixed in OpenSSL 1.1.1c (Affected 1.1.1-1.1.1b). Fixed in OpenSSL 1.1.0k (Affected 1.1.0-1.1.0j).",
        "Severity": "HIGH",
        "CweIDs": [
          "CWE-310"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:M/Au:N/C:P/I:P/A:N",
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:N",
            "V2Score": 5.8,
            "V3Score": 7.4
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:L/AC:H/PR:N/UI:N/S:U/C:N/I:L/A:N",
            "V3Score": 2.9
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-07/msg00056.html",
          "https://access.redhat.com/errata/RHSA-2019:3700",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1543",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=ee22257b1418438ebaf54df98af4e24f494d1809",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=f426625b6ae9a7831010750490a5f0ad689c5ba3",
          "https://linux.oracle.com/cve/CVE-2019-1543.html",
          "https://linux.oracle.com/errata/ELSA-2019-3700.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/Y3IVFGSERAZLNJCK35TEM2R4726XIH3Z/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/ZBEV5QGDRFUZDMNECFXUSN5FMYOZDE4V/",
          "https://seclists.org/bugtraq/2019/Jul/3",
          "https://www.debian.org/security/2019/dsa-4475",
          "https://www.openssl.org/news/secadv/20190306.txt",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpujul2019-5072835.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html"
        ],
        "PublishedDate": "2019-03-06T21:29:00Z",
        "LastModifiedDate": "2019-06-03T20:29:00Z"
      },
      {
        "VulnerabilityID": "CVE-2020-1967",
        "PkgName": "libcrypto1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1g-r0",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2020-1967",
        "Title": "openssl: Segmentation fault in SSL_check_chain causes denial of service",
        "Description": "Server or client applications that call the SSL_check_chain() function during or after a TLS 1.3 handshake may crash due to a NULL pointer dereference as a result of incorrect handling of the \"signature_algorithms_cert\" TLS extension. The crash occurs if an invalid or unrecognised signature algorithm is received from the peer. This could be exploited by a malicious peer in a Denial of Service attack. OpenSSL version 1.1.1d, 1.1.1e, and 1.1.1f are affected by this issue. This issue did not affect OpenSSL versions prior to 1.1.1d. Fixed in OpenSSL 1.1.1g (Affected 1.1.1d-1.1.1f).",
        "Severity": "HIGH",
        "CweIDs": [
          "CWE-476"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 5,
            "V3Score": 7.5
          },
          "redhat": {
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V3Score": 7.5
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2020-07/msg00004.html",
          "http://lists.opensuse.org/opensuse-security-announce/2020-07/msg00011.html",
          "http://packetstormsecurity.com/files/157527/OpenSSL-signature_algorithms_cert-Denial-Of-Service.html",
          "http://seclists.org/fulldisclosure/2020/May/5",
          "http://www.openwall.com/lists/oss-security/2020/04/22/2",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-1967",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=eb563247aef3e83dda7679c43f9649270462e5b1",
          "https://github.com/irsl/CVE-2020-1967",
          "https://kb.pulsesecure.net/articles/Pulse_Security_Advisories/SA44440",
          "https://lists.apache.org/thread.html/r66ea9c436da150683432db5fbc8beb8ae01886c6459ac30c2cea7345@%3Cdev.tomcat.apache.org%3E",
          "https://lists.apache.org/thread.html/r94d6ac3f010a38fccf4f432b12180a13fa1cf303559bd805648c9064@%3Cdev.tomcat.apache.org%3E",
          "https://lists.apache.org/thread.html/r9a41e304992ce6aec6585a87842b4f2e692604f5c892c37e3b0587ee@%3Cdev.tomcat.apache.org%3E",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/DDHOAATPWJCXRNFMJ2SASDBBNU5RJONY/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/EXDDAOWSAIEFQNBHWYE6PPYFV4QXGMCD/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/XVEP3LAK4JSPRXFO4QF4GG2IVXADV3SO/",
          "https://security.FreeBSD.org/advisories/FreeBSD-SA-20:11.openssl.asc",
          "https://security.gentoo.org/glsa/202004-10",
          "https://security.netapp.com/advisory/ntap-20200424-0003/",
          "https://security.netapp.com/advisory/ntap-20200717-0004/",
          "https://www.debian.org/security/2020/dsa-4661",
          "https://www.openssl.org/news/secadv/20200421.txt",
          "https://www.oracle.com/security-alerts/cpujan2021.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www.synology.com/security/advisory/Synology_SA_20_05",
          "https://www.synology.com/security/advisory/Synology_SA_20_05_OpenSSL",
          "https://www.tenable.com/security/tns-2020-03",
          "https://www.tenable.com/security/tns-2020-04",
          "https://www.tenable.com/security/tns-2020-11"
        ],
        "PublishedDate": "2020-04-21T14:15:00Z",
        "LastModifiedDate": "2021-02-09T15:14:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1547",
        "PkgName": "libcrypto1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1d-r0",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1547",
        "Title": "openssl: side-channel weak encryption vulnerability",
        "Description": "Normally in OpenSSL EC groups always have a co-factor present and this is used in side channel resistant code paths. However, in some cases, it is possible to construct a group using explicit parameters (instead of using a named curve). In those cases it is possible that such a group does not have the cofactor present. This can occur even where all the parameters match a known named curve. If such a curve is used then OpenSSL falls back to non-side channel resistant code paths which may result in full key recovery during an ECDSA signature operation. In order to be vulnerable an attacker would have to have the ability to time the creation of a large number of signatures where explicit parameters with no co-factor present are in use by an application using libcrypto. For the avoidance of doubt libssl is not vulnerable because explicit parameters are never used. Fixed in OpenSSL 1.1.1d (Affected 1.1.1-1.1.1c). Fixed in OpenSSL 1.1.0l (Affected 1.1.0-1.1.0k). Fixed in OpenSSL 1.0.2t (Affected 1.0.2-1.0.2s).",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-311"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:L/AC:M/Au:N/C:P/I:N/A:N",
            "V3Vector": "CVSS:3.1/AV:L/AC:H/PR:L/UI:N/S:U/C:H/I:N/A:N",
            "V2Score": 1.9,
            "V3Score": 4.7
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:L/AC:L/PR:L/UI:N/S:U/C:H/I:N/A:N",
            "V3Score": 5.5
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00054.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00072.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00012.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00016.html",
          "http://packetstormsecurity.com/files/154467/Slackware-Security-Advisory-openssl-Updates.html",
          "https://arxiv.org/abs/1909.01785",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1547",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=21c856b75d81eff61aa63b4f036bb64a85bf6d46",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=30c22fa8b1d840036b8e203585738df62a03cec8",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=7c1709c2da5414f5b6133d00a03fc8c5bf996c7a",
          "https://linux.oracle.com/cve/CVE-2019-1547.html",
          "https://linux.oracle.com/errata/ELSA-2020-1840.html",
          "https://lists.debian.org/debian-lts-announce/2019/09/msg00026.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/GY6SNRJP2S7Y42GIIDO3HXPNMDYN2U3A/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/ZN4VVQJ3JDCHGIHV4Y2YTXBYQZ6PWQ7E/",
          "https://seclists.org/bugtraq/2019/Oct/0",
          "https://seclists.org/bugtraq/2019/Oct/1",
          "https://seclists.org/bugtraq/2019/Sep/25",
          "https://security.gentoo.org/glsa/201911-04",
          "https://security.netapp.com/advisory/ntap-20190919-0002/",
          "https://security.netapp.com/advisory/ntap-20200122-0002/",
          "https://security.netapp.com/advisory/ntap-20200416-0003/",
          "https://support.f5.com/csp/article/K73422160?utm_source=f5support\u0026amp;utm_medium=RSS",
          "https://usn.ubuntu.com/4376-1/",
          "https://usn.ubuntu.com/4376-2/",
          "https://usn.ubuntu.com/4504-1/",
          "https://usn.ubuntu.com/usn/usn-4376-1",
          "https://usn.ubuntu.com/usn/usn-4376-2",
          "https://usn.ubuntu.com/usn/usn-4504-1",
          "https://www.debian.org/security/2019/dsa-4539",
          "https://www.debian.org/security/2019/dsa-4540",
          "https://www.openssl.org/news/secadv/20190910.txt",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpujan2020.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html",
          "https://www.tenable.com/security/tns-2019-08",
          "https://www.tenable.com/security/tns-2019-09"
        ],
        "PublishedDate": "2019-09-10T17:15:00Z",
        "LastModifiedDate": "2020-10-20T22:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1549",
        "PkgName": "libcrypto1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1d-r0",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1549",
        "Title": "openssl: information disclosure in fork()",
        "Description": "OpenSSL 1.1.1 introduced a rewritten random number generator (RNG). This was intended to include protection in the event of a fork() system call in order to ensure that the parent and child processes did not share the same RNG state. However this protection was not being used in the default case. A partial mitigation for this issue is that the output from a high precision timer is mixed into the RNG state so the likelihood of a parent and child process sharing state is significantly reduced. If an application already calls OPENSSL_init_crypto() explicitly using OPENSSL_INIT_ATFORK then this problem does not occur at all. Fixed in OpenSSL 1.1.1d (Affected 1.1.1-1.1.1c).",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-330"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:P/I:N/A:N",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:L/I:N/A:N",
            "V2Score": 5,
            "V3Score": 5.3
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:L/A:N",
            "V3Score": 4.8
          }
        },
        "References": [
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1549",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=1b0fe00e2704b5e20334a16d3c9099d1ba2ef1be",
          "https://linux.oracle.com/cve/CVE-2019-1549.html",
          "https://linux.oracle.com/errata/ELSA-2020-1840.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/GY6SNRJP2S7Y42GIIDO3HXPNMDYN2U3A/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/ZN4VVQJ3JDCHGIHV4Y2YTXBYQZ6PWQ7E/",
          "https://seclists.org/bugtraq/2019/Oct/1",
          "https://security.netapp.com/advisory/ntap-20190919-0002/",
          "https://support.f5.com/csp/article/K44070243",
          "https://support.f5.com/csp/article/K44070243?utm_source=f5support\u0026amp;utm_medium=RSS",
          "https://usn.ubuntu.com/4376-1/",
          "https://usn.ubuntu.com/usn/usn-4376-1",
          "https://www.debian.org/security/2019/dsa-4539",
          "https://www.openssl.org/news/secadv/20190910.txt",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpujan2020.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html"
        ],
        "PublishedDate": "2019-09-10T17:15:00Z",
        "LastModifiedDate": "2020-10-20T22:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1551",
        "PkgName": "libcrypto1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1d-r2",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1551",
        "Title": "openssl: Integer overflow in RSAZ modular exponentiation on x86_64",
        "Description": "There is an overflow bug in the x64_64 Montgomery squaring procedure used in exponentiation with 512-bit moduli. No EC algorithms are affected. Analysis suggests that attacks against 2-prime RSA1024, 3-prime RSA1536, and DSA1024 as a result of this defect would be very difficult to perform and are not believed likely. Attacks against DH512 are considered just feasible. However, for an attack the target would have to re-use the DH512 private key, which is not recommended anyway. Also applications directly using the low level API BN_mod_exp may be affected if they use BN_FLG_CONSTTIME. Fixed in OpenSSL 1.1.1e (Affected 1.1.1-1.1.1d). Fixed in OpenSSL 1.0.2u (Affected 1.0.2-1.0.2t).",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-200"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:P/I:N/A:N",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:L/I:N/A:N",
            "V2Score": 5,
            "V3Score": 5.3
          },
          "redhat": {
            "V3Vector": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:L/A:N",
            "V3Score": 4.8
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2020-01/msg00030.html",
          "http://packetstormsecurity.com/files/155754/Slackware-Security-Advisory-openssl-Updates.html",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1551",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=419102400a2811582a7a3d4a4e317d72e5ce0a8f",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=f1c5eea8a817075d31e43f5876993c6710238c98",
          "https://github.com/openssl/openssl/pull/10575",
          "https://linux.oracle.com/cve/CVE-2019-1551.html",
          "https://linux.oracle.com/errata/ELSA-2020-4514.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/DDHOAATPWJCXRNFMJ2SASDBBNU5RJONY/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/EXDDAOWSAIEFQNBHWYE6PPYFV4QXGMCD/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/XVEP3LAK4JSPRXFO4QF4GG2IVXADV3SO/",
          "https://seclists.org/bugtraq/2019/Dec/39",
          "https://seclists.org/bugtraq/2019/Dec/46",
          "https://security.gentoo.org/glsa/202004-10",
          "https://security.netapp.com/advisory/ntap-20191210-0001/",
          "https://usn.ubuntu.com/4376-1/",
          "https://usn.ubuntu.com/4504-1/",
          "https://usn.ubuntu.com/usn/usn-4376-1",
          "https://usn.ubuntu.com/usn/usn-4504-1",
          "https://www.debian.org/security/2019/dsa-4594",
          "https://www.debian.org/security/2021/dsa-4855",
          "https://www.openssl.org/news/secadv/20191206.txt",
          "https://www.oracle.com/security-alerts/cpujan2021.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.tenable.com/security/tns-2019-09",
          "https://www.tenable.com/security/tns-2020-03",
          "https://www.tenable.com/security/tns-2020-11"
        ],
        "PublishedDate": "2019-12-06T18:15:00Z",
        "LastModifiedDate": "2021-02-24T00:44:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1563",
        "PkgName": "libcrypto1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1d-r0",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1563",
        "Title": "openssl: information disclosure in PKCS7_dataDecode and CMS_decrypt_set1_pkey",
        "Description": "In situations where an attacker receives automated notification of the success or failure of a decryption attempt an attacker, after sending a very large number of messages to be decrypted, can recover a CMS/PKCS7 transported encryption key or decrypt any RSA encrypted message that was encrypted with the public RSA key, using a Bleichenbacher padding oracle attack. Applications are not affected if they use a certificate together with the private RSA key to the CMS_decrypt or PKCS7_decrypt functions to select the correct recipient info to decrypt. Fixed in OpenSSL 1.1.1d (Affected 1.1.1-1.1.1c). Fixed in OpenSSL 1.1.0l (Affected 1.1.0-1.1.0k). Fixed in OpenSSL 1.0.2t (Affected 1.0.2-1.0.2s).",
        "Severity": "LOW",
        "CweIDs": [
          "CWE-311"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:M/Au:N/C:P/I:N/A:N",
            "V3Vector": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:N/A:N",
            "V2Score": 4.3,
            "V3Score": 3.7
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:N/A:N",
            "V3Score": 3.7
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00054.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00072.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00012.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00016.html",
          "http://packetstormsecurity.com/files/154467/Slackware-Security-Advisory-openssl-Updates.html",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1563",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=08229ad838c50f644d7e928e2eef147b4308ad64",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=631f94db0065c78181ca9ba5546ebc8bb3884b97",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=e21f8cf78a125cd3c8c0d1a1a6c8bb0b901f893f",
          "https://linux.oracle.com/cve/CVE-2019-1563.html",
          "https://linux.oracle.com/errata/ELSA-2020-1840.html",
          "https://lists.debian.org/debian-lts-announce/2019/09/msg00026.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/GY6SNRJP2S7Y42GIIDO3HXPNMDYN2U3A/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/ZN4VVQJ3JDCHGIHV4Y2YTXBYQZ6PWQ7E/",
          "https://seclists.org/bugtraq/2019/Oct/0",
          "https://seclists.org/bugtraq/2019/Oct/1",
          "https://seclists.org/bugtraq/2019/Sep/25",
          "https://security.gentoo.org/glsa/201911-04",
          "https://security.netapp.com/advisory/ntap-20190919-0002/",
          "https://support.f5.com/csp/article/K97324400?utm_source=f5support\u0026amp;utm_medium=RSS",
          "https://usn.ubuntu.com/4376-1/",
          "https://usn.ubuntu.com/4376-2/",
          "https://usn.ubuntu.com/4504-1/",
          "https://usn.ubuntu.com/usn/usn-4376-1",
          "https://usn.ubuntu.com/usn/usn-4376-2",
          "https://usn.ubuntu.com/usn/usn-4504-1",
          "https://www.debian.org/security/2019/dsa-4539",
          "https://www.debian.org/security/2019/dsa-4540",
          "https://www.openssl.org/news/secadv/20190910.txt",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpujan2020.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html",
          "https://www.tenable.com/security/tns-2019-09"
        ],
        "PublishedDate": "2019-09-10T17:15:00Z",
        "LastModifiedDate": "2020-10-20T22:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1543",
        "PkgName": "libssl1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1b-r1",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1543",
        "Title": "openssl: ChaCha20-Poly1305 with long nonces",
        "Description": "ChaCha20-Poly1305 is an AEAD cipher, and requires a unique nonce input for every encryption operation. RFC 7539 specifies that the nonce value (IV) should be 96 bits (12 bytes). OpenSSL allows a variable nonce length and front pads the nonce with 0 bytes if it is less than 12 bytes. However it also incorrectly allows a nonce to be set of up to 16 bytes. In this case only the last 12 bytes are significant and any additional leading bytes are ignored. It is a requirement of using this cipher that nonce values are unique. Messages encrypted using a reused nonce value are susceptible to serious confidentiality and integrity attacks. If an application changes the default nonce length to be longer than 12 bytes and then makes a change to the leading bytes of the nonce expecting the new value to be a new unique nonce then such an application could inadvertently encrypt messages with a reused nonce. Additionally the ignored bytes in a long nonce are not covered by the integrity guarantee of this cipher. Any application that relies on the integrity of these ignored leading bytes of a long nonce may be further affected. Any OpenSSL internal use of this cipher, including in SSL/TLS, is safe because no such use sets such a long nonce value. However user applications that use this cipher directly and set a non-default nonce length to be longer than 12 bytes may be vulnerable. OpenSSL versions 1.1.1 and 1.1.0 are affected by this issue. Due to the limited scope of affected deployments this has been assessed as low severity and therefore we are not creating new releases at this time. Fixed in OpenSSL 1.1.1c (Affected 1.1.1-1.1.1b). Fixed in OpenSSL 1.1.0k (Affected 1.1.0-1.1.0j).",
        "Severity": "HIGH",
        "CweIDs": [
          "CWE-310"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:M/Au:N/C:P/I:P/A:N",
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:N",
            "V2Score": 5.8,
            "V3Score": 7.4
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:L/AC:H/PR:N/UI:N/S:U/C:N/I:L/A:N",
            "V3Score": 2.9
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-07/msg00056.html",
          "https://access.redhat.com/errata/RHSA-2019:3700",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1543",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=ee22257b1418438ebaf54df98af4e24f494d1809",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=f426625b6ae9a7831010750490a5f0ad689c5ba3",
          "https://linux.oracle.com/cve/CVE-2019-1543.html",
          "https://linux.oracle.com/errata/ELSA-2019-3700.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/Y3IVFGSERAZLNJCK35TEM2R4726XIH3Z/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/ZBEV5QGDRFUZDMNECFXUSN5FMYOZDE4V/",
          "https://seclists.org/bugtraq/2019/Jul/3",
          "https://www.debian.org/security/2019/dsa-4475",
          "https://www.openssl.org/news/secadv/20190306.txt",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpujul2019-5072835.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html"
        ],
        "PublishedDate": "2019-03-06T21:29:00Z",
        "LastModifiedDate": "2019-06-03T20:29:00Z"
      },
      {
        "VulnerabilityID": "CVE-2020-1967",
        "PkgName": "libssl1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1g-r0",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2020-1967",
        "Title": "openssl: Segmentation fault in SSL_check_chain causes denial of service",
        "Description": "Server or client applications that call the SSL_check_chain() function during or after a TLS 1.3 handshake may crash due to a NULL pointer dereference as a result of incorrect handling of the \"signature_algorithms_cert\" TLS extension. The crash occurs if an invalid or unrecognised signature algorithm is received from the peer. This could be exploited by a malicious peer in a Denial of Service attack. OpenSSL version 1.1.1d, 1.1.1e, and 1.1.1f are affected by this issue. This issue did not affect OpenSSL versions prior to 1.1.1d. Fixed in OpenSSL 1.1.1g (Affected 1.1.1d-1.1.1f).",
        "Severity": "HIGH",
        "CweIDs": [
          "CWE-476"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 5,
            "V3Score": 7.5
          },
          "redhat": {
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V3Score": 7.5
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2020-07/msg00004.html",
          "http://lists.opensuse.org/opensuse-security-announce/2020-07/msg00011.html",
          "http://packetstormsecurity.com/files/157527/OpenSSL-signature_algorithms_cert-Denial-Of-Service.html",
          "http://seclists.org/fulldisclosure/2020/May/5",
          "http://www.openwall.com/lists/oss-security/2020/04/22/2",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-1967",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=eb563247aef3e83dda7679c43f9649270462e5b1",
          "https://github.com/irsl/CVE-2020-1967",
          "https://kb.pulsesecure.net/articles/Pulse_Security_Advisories/SA44440",
          "https://lists.apache.org/thread.html/r66ea9c436da150683432db5fbc8beb8ae01886c6459ac30c2cea7345@%3Cdev.tomcat.apache.org%3E",
          "https://lists.apache.org/thread.html/r94d6ac3f010a38fccf4f432b12180a13fa1cf303559bd805648c9064@%3Cdev.tomcat.apache.org%3E",
          "https://lists.apache.org/thread.html/r9a41e304992ce6aec6585a87842b4f2e692604f5c892c37e3b0587ee@%3Cdev.tomcat.apache.org%3E",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/DDHOAATPWJCXRNFMJ2SASDBBNU5RJONY/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/EXDDAOWSAIEFQNBHWYE6PPYFV4QXGMCD/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/XVEP3LAK4JSPRXFO4QF4GG2IVXADV3SO/",
          "https://security.FreeBSD.org/advisories/FreeBSD-SA-20:11.openssl.asc",
          "https://security.gentoo.org/glsa/202004-10",
          "https://security.netapp.com/advisory/ntap-20200424-0003/",
          "https://security.netapp.com/advisory/ntap-20200717-0004/",
          "https://www.debian.org/security/2020/dsa-4661",
          "https://www.openssl.org/news/secadv/20200421.txt",
          "https://www.oracle.com/security-alerts/cpujan2021.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www.synology.com/security/advisory/Synology_SA_20_05",
          "https://www.synology.com/security/advisory/Synology_SA_20_05_OpenSSL",
          "https://www.tenable.com/security/tns-2020-03",
          "https://www.tenable.com/security/tns-2020-04",
          "https://www.tenable.com/security/tns-2020-11"
        ],
        "PublishedDate": "2020-04-21T14:15:00Z",
        "LastModifiedDate": "2021-02-09T15:14:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1547",
        "PkgName": "libssl1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1d-r0",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1547",
        "Title": "openssl: side-channel weak encryption vulnerability",
        "Description": "Normally in OpenSSL EC groups always have a co-factor present and this is used in side channel resistant code paths. However, in some cases, it is possible to construct a group using explicit parameters (instead of using a named curve). In those cases it is possible that such a group does not have the cofactor present. This can occur even where all the parameters match a known named curve. If such a curve is used then OpenSSL falls back to non-side channel resistant code paths which may result in full key recovery during an ECDSA signature operation. In order to be vulnerable an attacker would have to have the ability to time the creation of a large number of signatures where explicit parameters with no co-factor present are in use by an application using libcrypto. For the avoidance of doubt libssl is not vulnerable because explicit parameters are never used. Fixed in OpenSSL 1.1.1d (Affected 1.1.1-1.1.1c). Fixed in OpenSSL 1.1.0l (Affected 1.1.0-1.1.0k). Fixed in OpenSSL 1.0.2t (Affected 1.0.2-1.0.2s).",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-311"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:L/AC:M/Au:N/C:P/I:N/A:N",
            "V3Vector": "CVSS:3.1/AV:L/AC:H/PR:L/UI:N/S:U/C:H/I:N/A:N",
            "V2Score": 1.9,
            "V3Score": 4.7
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:L/AC:L/PR:L/UI:N/S:U/C:H/I:N/A:N",
            "V3Score": 5.5
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00054.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00072.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00012.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00016.html",
          "http://packetstormsecurity.com/files/154467/Slackware-Security-Advisory-openssl-Updates.html",
          "https://arxiv.org/abs/1909.01785",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1547",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=21c856b75d81eff61aa63b4f036bb64a85bf6d46",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=30c22fa8b1d840036b8e203585738df62a03cec8",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=7c1709c2da5414f5b6133d00a03fc8c5bf996c7a",
          "https://linux.oracle.com/cve/CVE-2019-1547.html",
          "https://linux.oracle.com/errata/ELSA-2020-1840.html",
          "https://lists.debian.org/debian-lts-announce/2019/09/msg00026.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/GY6SNRJP2S7Y42GIIDO3HXPNMDYN2U3A/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/ZN4VVQJ3JDCHGIHV4Y2YTXBYQZ6PWQ7E/",
          "https://seclists.org/bugtraq/2019/Oct/0",
          "https://seclists.org/bugtraq/2019/Oct/1",
          "https://seclists.org/bugtraq/2019/Sep/25",
          "https://security.gentoo.org/glsa/201911-04",
          "https://security.netapp.com/advisory/ntap-20190919-0002/",
          "https://security.netapp.com/advisory/ntap-20200122-0002/",
          "https://security.netapp.com/advisory/ntap-20200416-0003/",
          "https://support.f5.com/csp/article/K73422160?utm_source=f5support\u0026amp;utm_medium=RSS",
          "https://usn.ubuntu.com/4376-1/",
          "https://usn.ubuntu.com/4376-2/",
          "https://usn.ubuntu.com/4504-1/",
          "https://usn.ubuntu.com/usn/usn-4376-1",
          "https://usn.ubuntu.com/usn/usn-4376-2",
          "https://usn.ubuntu.com/usn/usn-4504-1",
          "https://www.debian.org/security/2019/dsa-4539",
          "https://www.debian.org/security/2019/dsa-4540",
          "https://www.openssl.org/news/secadv/20190910.txt",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpujan2020.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html",
          "https://www.tenable.com/security/tns-2019-08",
          "https://www.tenable.com/security/tns-2019-09"
        ],
        "PublishedDate": "2019-09-10T17:15:00Z",
        "LastModifiedDate": "2020-10-20T22:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1549",
        "PkgName": "libssl1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1d-r0",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1549",
        "Title": "openssl: information disclosure in fork()",
        "Description": "OpenSSL 1.1.1 introduced a rewritten random number generator (RNG). This was intended to include protection in the event of a fork() system call in order to ensure that the parent and child processes did not share the same RNG state. However this protection was not being used in the default case. A partial mitigation for this issue is that the output from a high precision timer is mixed into the RNG state so the likelihood of a parent and child process sharing state is significantly reduced. If an application already calls OPENSSL_init_crypto() explicitly using OPENSSL_INIT_ATFORK then this problem does not occur at all. Fixed in OpenSSL 1.1.1d (Affected 1.1.1-1.1.1c).",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-330"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:P/I:N/A:N",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:L/I:N/A:N",
            "V2Score": 5,
            "V3Score": 5.3
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:L/A:N",
            "V3Score": 4.8
          }
        },
        "References": [
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1549",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=1b0fe00e2704b5e20334a16d3c9099d1ba2ef1be",
          "https://linux.oracle.com/cve/CVE-2019-1549.html",
          "https://linux.oracle.com/errata/ELSA-2020-1840.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/GY6SNRJP2S7Y42GIIDO3HXPNMDYN2U3A/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/ZN4VVQJ3JDCHGIHV4Y2YTXBYQZ6PWQ7E/",
          "https://seclists.org/bugtraq/2019/Oct/1",
          "https://security.netapp.com/advisory/ntap-20190919-0002/",
          "https://support.f5.com/csp/article/K44070243",
          "https://support.f5.com/csp/article/K44070243?utm_source=f5support\u0026amp;utm_medium=RSS",
          "https://usn.ubuntu.com/4376-1/",
          "https://usn.ubuntu.com/usn/usn-4376-1",
          "https://www.debian.org/security/2019/dsa-4539",
          "https://www.openssl.org/news/secadv/20190910.txt",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpujan2020.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html"
        ],
        "PublishedDate": "2019-09-10T17:15:00Z",
        "LastModifiedDate": "2020-10-20T22:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1551",
        "PkgName": "libssl1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1d-r2",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1551",
        "Title": "openssl: Integer overflow in RSAZ modular exponentiation on x86_64",
        "Description": "There is an overflow bug in the x64_64 Montgomery squaring procedure used in exponentiation with 512-bit moduli. No EC algorithms are affected. Analysis suggests that attacks against 2-prime RSA1024, 3-prime RSA1536, and DSA1024 as a result of this defect would be very difficult to perform and are not believed likely. Attacks against DH512 are considered just feasible. However, for an attack the target would have to re-use the DH512 private key, which is not recommended anyway. Also applications directly using the low level API BN_mod_exp may be affected if they use BN_FLG_CONSTTIME. Fixed in OpenSSL 1.1.1e (Affected 1.1.1-1.1.1d). Fixed in OpenSSL 1.0.2u (Affected 1.0.2-1.0.2t).",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-200"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:P/I:N/A:N",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:L/I:N/A:N",
            "V2Score": 5,
            "V3Score": 5.3
          },
          "redhat": {
            "V3Vector": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:L/A:N",
            "V3Score": 4.8
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2020-01/msg00030.html",
          "http://packetstormsecurity.com/files/155754/Slackware-Security-Advisory-openssl-Updates.html",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1551",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=419102400a2811582a7a3d4a4e317d72e5ce0a8f",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=f1c5eea8a817075d31e43f5876993c6710238c98",
          "https://github.com/openssl/openssl/pull/10575",
          "https://linux.oracle.com/cve/CVE-2019-1551.html",
          "https://linux.oracle.com/errata/ELSA-2020-4514.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/DDHOAATPWJCXRNFMJ2SASDBBNU5RJONY/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/EXDDAOWSAIEFQNBHWYE6PPYFV4QXGMCD/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/XVEP3LAK4JSPRXFO4QF4GG2IVXADV3SO/",
          "https://seclists.org/bugtraq/2019/Dec/39",
          "https://seclists.org/bugtraq/2019/Dec/46",
          "https://security.gentoo.org/glsa/202004-10",
          "https://security.netapp.com/advisory/ntap-20191210-0001/",
          "https://usn.ubuntu.com/4376-1/",
          "https://usn.ubuntu.com/4504-1/",
          "https://usn.ubuntu.com/usn/usn-4376-1",
          "https://usn.ubuntu.com/usn/usn-4504-1",
          "https://www.debian.org/security/2019/dsa-4594",
          "https://www.debian.org/security/2021/dsa-4855",
          "https://www.openssl.org/news/secadv/20191206.txt",
          "https://www.oracle.com/security-alerts/cpujan2021.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.tenable.com/security/tns-2019-09",
          "https://www.tenable.com/security/tns-2020-03",
          "https://www.tenable.com/security/tns-2020-11"
        ],
        "PublishedDate": "2019-12-06T18:15:00Z",
        "LastModifiedDate": "2021-02-24T00:44:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-1563",
        "PkgName": "libssl1.1",
        "InstalledVersion": "1.1.1a-r1",
        "FixedVersion": "1.1.1d-r0",
        "Layer": {
          "DiffID": "sha256:bcf2f368fe234217249e00ad9d762d8f1a3156d60c442ed92079fa5b120634a1"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-1563",
        "Title": "openssl: information disclosure in PKCS7_dataDecode and CMS_decrypt_set1_pkey",
        "Description": "In situations where an attacker receives automated notification of the success or failure of a decryption attempt an attacker, after sending a very large number of messages to be decrypted, can recover a CMS/PKCS7 transported encryption key or decrypt any RSA encrypted message that was encrypted with the public RSA key, using a Bleichenbacher padding oracle attack. Applications are not affected if they use a certificate together with the private RSA key to the CMS_decrypt or PKCS7_decrypt functions to select the correct recipient info to decrypt. Fixed in OpenSSL 1.1.1d (Affected 1.1.1-1.1.1c). Fixed in OpenSSL 1.1.0l (Affected 1.1.0-1.1.0k). Fixed in OpenSSL 1.0.2t (Affected 1.0.2-1.0.2s).",
        "Severity": "LOW",
        "CweIDs": [
          "CWE-311"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:M/Au:N/C:P/I:N/A:N",
            "V3Vector": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:N/A:N",
            "V2Score": 4.3,
            "V3Score": 3.7
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:N/A:N",
            "V3Score": 3.7
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00054.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-09/msg00072.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00012.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00016.html",
          "http://packetstormsecurity.com/files/154467/Slackware-Security-Advisory-openssl-Updates.html",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-1563",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=08229ad838c50f644d7e928e2eef147b4308ad64",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=631f94db0065c78181ca9ba5546ebc8bb3884b97",
          "https://git.openssl.org/gitweb/?p=openssl.git;a=commitdiff;h=e21f8cf78a125cd3c8c0d1a1a6c8bb0b901f893f",
          "https://linux.oracle.com/cve/CVE-2019-1563.html",
          "https://linux.oracle.com/errata/ELSA-2020-1840.html",
          "https://lists.debian.org/debian-lts-announce/2019/09/msg00026.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/GY6SNRJP2S7Y42GIIDO3HXPNMDYN2U3A/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/ZN4VVQJ3JDCHGIHV4Y2YTXBYQZ6PWQ7E/",
          "https://seclists.org/bugtraq/2019/Oct/0",
          "https://seclists.org/bugtraq/2019/Oct/1",
          "https://seclists.org/bugtraq/2019/Sep/25",
          "https://security.gentoo.org/glsa/201911-04",
          "https://security.netapp.com/advisory/ntap-20190919-0002/",
          "https://support.f5.com/csp/article/K97324400?utm_source=f5support\u0026amp;utm_medium=RSS",
          "https://usn.ubuntu.com/4376-1/",
          "https://usn.ubuntu.com/4376-2/",
          "https://usn.ubuntu.com/4504-1/",
          "https://usn.ubuntu.com/usn/usn-4376-1",
          "https://usn.ubuntu.com/usn/usn-4376-2",
          "https://usn.ubuntu.com/usn/usn-4504-1",
          "https://www.debian.org/security/2019/dsa-4539",
          "https://www.debian.org/security/2019/dsa-4540",
          "https://www.openssl.org/news/secadv/20190910.txt",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpujan2020.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html",
          "https://www.tenable.com/security/tns-2019-09"
        ],
        "PublishedDate": "2019-09-10T17:15:00Z",
        "LastModifiedDate": "2020-10-20T22:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-14697",
        "PkgName": "musl",
        "InstalledVersion": "1.1.20-r4",
        "FixedVersion": "1.1.20-r5",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-14697",
        "Description": "musl libc through 1.1.23 has an x87 floating-point stack adjustment imbalance, related to the math/i386/ directory. In some cases, use of this library could introduce out-of-bounds writes that are not present in an application's source code.",
        "Severity": "CRITICAL",
        "CweIDs": [
          "CWE-787"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:P/I:P/A:P",
            "V3Vector": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
            "V2Score": 7.5,
            "V3Score": 9.8
          }
        },
        "References": [
          "http://www.openwall.com/lists/oss-security/2019/08/06/4",
          "https://security.gentoo.org/glsa/202003-13",
          "https://www.openwall.com/lists/musl/2019/08/06/1"
        ],
        "PublishedDate": "2019-08-06T16:15:00Z",
        "LastModifiedDate": "2020-03-14T19:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2020-28928",
        "PkgName": "musl",
        "InstalledVersion": "1.1.20-r4",
        "FixedVersion": "1.1.20-r6",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2020-28928",
        "Description": "In musl libc through 1.2.1, wcsnrtombs mishandles particular combinations of destination buffer size and source character limit, as demonstrated by an invalid write access (buffer overflow).",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-787"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:L/AC:L/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 2.1,
            "V3Score": 5.5
          }
        },
        "References": [
          "http://www.openwall.com/lists/oss-security/2020/11/20/4",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-28928",
          "https://lists.debian.org/debian-lts-announce/2020/11/msg00050.html",
          "https://musl.libc.org/releases.html"
        ],
        "PublishedDate": "2020-11-24T18:15:00Z",
        "LastModifiedDate": "2021-03-11T18:16:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-14697",
        "PkgName": "musl-utils",
        "InstalledVersion": "1.1.20-r4",
        "FixedVersion": "1.1.20-r5",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-14697",
        "Description": "musl libc through 1.1.23 has an x87 floating-point stack adjustment imbalance, related to the math/i386/ directory. In some cases, use of this library could introduce out-of-bounds writes that are not present in an application's source code.",
        "Severity": "CRITICAL",
        "CweIDs": [
          "CWE-787"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:P/I:P/A:P",
            "V3Vector": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
            "V2Score": 7.5,
            "V3Score": 9.8
          }
        },
        "References": [
          "http://www.openwall.com/lists/oss-security/2019/08/06/4",
          "https://security.gentoo.org/glsa/202003-13",
          "https://www.openwall.com/lists/musl/2019/08/06/1"
        ],
        "PublishedDate": "2019-08-06T16:15:00Z",
        "LastModifiedDate": "2020-03-14T19:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2020-28928",
        "PkgName": "musl-utils",
        "InstalledVersion": "1.1.20-r4",
        "FixedVersion": "1.1.20-r6",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2020-28928",
        "Description": "In musl libc through 1.2.1, wcsnrtombs mishandles particular combinations of destination buffer size and source character limit, as demonstrated by an invalid write access (buffer overflow).",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-787"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:L/AC:L/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 2.1,
            "V3Score": 5.5
          }
        },
        "References": [
          "http://www.openwall.com/lists/oss-security/2020/11/20/4",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-28928",
          "https://lists.debian.org/debian-lts-announce/2020/11/msg00050.html",
          "https://musl.libc.org/releases.html"
        ],
        "PublishedDate": "2020-11-24T18:15:00Z",
        "LastModifiedDate": "2021-03-11T18:16:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-8457",
        "PkgName": "sqlite-libs",
        "InstalledVersion": "3.26.0-r3",
        "FixedVersion": "3.28.0-r0",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-8457",
        "Title": "sqlite: heap out-of-bound read in function rtreenode()",
        "Description": "SQLite3 from 3.6.0 to and including 3.27.2 is vulnerable to heap out-of-bound read in the rtreenode() function when handling invalid rtree tables.",
        "Severity": "CRITICAL",
        "CweIDs": [
          "CWE-125"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:P/I:P/A:P",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
            "V2Score": 7.5,
            "V3Score": 9.8
          },
          "redhat": {
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V3Score": 7.5
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-06/msg00074.html",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-8457",
          "https://linux.oracle.com/cve/CVE-2019-8457.html",
          "https://linux.oracle.com/errata/ELSA-2020-1810.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/OPKYSWCOM3CL66RI76TYVIG6TJ263RXH/",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/SJPFGA45DI4F5MCF2OAACGH3HQOF4G3M/",
          "https://security.netapp.com/advisory/ntap-20190606-0002/",
          "https://usn.ubuntu.com/4004-1/",
          "https://usn.ubuntu.com/4004-2/",
          "https://usn.ubuntu.com/4019-1/",
          "https://usn.ubuntu.com/4019-2/",
          "https://usn.ubuntu.com/usn/usn-4004-1",
          "https://usn.ubuntu.com/usn/usn-4004-2",
          "https://usn.ubuntu.com/usn/usn-4019-1",
          "https://usn.ubuntu.com/usn/usn-4019-2",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpujan2020.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/technetwork/security-advisory/cpuoct2019-5072832.html",
          "https://www.sqlite.org/releaselog/3_28_0.html",
          "https://www.sqlite.org/src/info/90acdbfce9c08858"
        ],
        "PublishedDate": "2019-05-30T16:29:00Z",
        "LastModifiedDate": "2020-07-15T03:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-19244",
        "PkgName": "sqlite-libs",
        "InstalledVersion": "3.26.0-r3",
        "FixedVersion": "3.28.0-r2",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-19244",
        "Title": "sqlite: allows a crash if a sub-select uses both DISTINCT and window functions and also has certain ORDER BY usage",
        "Description": "sqlite3Select in select.c in SQLite 3.30.1 allows a crash if a sub-select uses both DISTINCT and window functions, and also has certain ORDER BY usage.",
        "Severity": "HIGH",
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 5,
            "V3Score": 7.5
          },
          "redhat": {
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V3Score": 7.5
          }
        },
        "References": [
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-19244",
          "https://github.com/sqlite/sqlite/commit/e59c562b3f6894f84c715772c4b116d7b5c01348",
          "https://usn.ubuntu.com/4205-1/",
          "https://usn.ubuntu.com/usn/usn-4205-1",
          "https://www.oracle.com/security-alerts/cpuapr2020.html"
        ],
        "PublishedDate": "2019-11-25T20:15:00Z",
        "LastModifiedDate": "2020-08-24T17:37:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-5018",
        "PkgName": "sqlite-libs",
        "InstalledVersion": "3.26.0-r3",
        "FixedVersion": "3.28.0-r0",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-5018",
        "Title": "sqlite: Use-after-free in window function leading to remote code execution",
        "Description": "An exploitable use after free vulnerability exists in the window function functionality of Sqlite3 3.26.0. A specially crafted SQL command can cause a use after free vulnerability, potentially resulting in remote code execution. An attacker can send a malicious SQL command to trigger this vulnerability.",
        "Severity": "HIGH",
        "CweIDs": [
          "CWE-416"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:M/Au:N/C:P/I:P/A:P",
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:H",
            "V2Score": 6.8,
            "V3Score": 8.1
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:H",
            "V3Score": 8.1
          }
        },
        "References": [
          "http://packetstormsecurity.com/files/152809/Sqlite3-Window-Function-Remote-Code-Execution.html",
          "http://www.securityfocus.com/bid/108294",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5018",
          "https://linux.oracle.com/cve/CVE-2019-5018.html",
          "https://linux.oracle.com/errata/ELSA-2020-4442.html",
          "https://security.gentoo.org/glsa/201908-09",
          "https://security.netapp.com/advisory/ntap-20190521-0001/",
          "https://talosintelligence.com/vulnerability_reports/TALOS-2019-0777",
          "https://usn.ubuntu.com/4205-1/",
          "https://usn.ubuntu.com/usn/usn-4205-1",
          "https://www.talosintelligence.com/vulnerability_reports/TALOS-2019-0777"
        ],
        "PublishedDate": "2019-05-10T19:29:00Z",
        "LastModifiedDate": "2019-05-21T19:29:00Z"
      },
      {
        "VulnerabilityID": "CVE-2020-11655",
        "PkgName": "sqlite-libs",
        "InstalledVersion": "3.26.0-r3",
        "FixedVersion": "3.28.0-r3",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2020-11655",
        "Title": "sqlite: malformed window-function query leads to DoS",
        "Description": "SQLite through 3.31.1 allows attackers to cause a denial of service (segmentation fault) via a malformed window-function query because the AggInfo object's initialization is mishandled.",
        "Severity": "HIGH",
        "CweIDs": [
          "CWE-665"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:L/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 5,
            "V3Score": 7.5
          },
          "redhat": {
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V3Score": 7.5
          }
        },
        "References": [
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-11655",
          "https://lists.debian.org/debian-lts-announce/2020/05/msg00006.html",
          "https://lists.debian.org/debian-lts-announce/2020/08/msg00037.html",
          "https://security.FreeBSD.org/advisories/FreeBSD-SA-20:22.sqlite.asc",
          "https://security.gentoo.org/glsa/202007-26",
          "https://security.netapp.com/advisory/ntap-20200416-0001/",
          "https://usn.ubuntu.com/4394-1/",
          "https://usn.ubuntu.com/usn/usn-4394-1",
          "https://www.oracle.com/security-alerts/cpujan2021.html",
          "https://www.oracle.com/security-alerts/cpujul2020.html",
          "https://www.oracle.com/security-alerts/cpuoct2020.html",
          "https://www3.sqlite.org/cgi/src/info/4a302b42c7bf5e11",
          "https://www3.sqlite.org/cgi/src/tktview?name=af4556bb5c"
        ],
        "PublishedDate": "2020-04-09T03:15:00Z",
        "LastModifiedDate": "2021-03-04T20:44:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-16168",
        "PkgName": "sqlite-libs",
        "InstalledVersion": "3.26.0-r3",
        "FixedVersion": "3.28.0-r1",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-16168",
        "Title": "sqlite: Division by zero in whereLoopAddBtreeIndex in sqlite3.c",
        "Description": "In SQLite through 3.29.0, whereLoopAddBtreeIndex in sqlite3.c can crash a browser or other application because of missing validation of a sqlite_stat1 sz field, aka a \"severe division by zero in the query planner.\"",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-369"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:M/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:R/S:U/C:N/I:N/A:H",
            "V2Score": 4.3,
            "V3Score": 6.5
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:N/AC:L/PR:N/UI:R/S:U/C:N/I:N/A:H",
            "V3Score": 6.5
          }
        },
        "References": [
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00032.html",
          "http://lists.opensuse.org/opensuse-security-announce/2019-10/msg00033.html",
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-16168",
          "https://linux.oracle.com/cve/CVE-2019-16168.html",
          "https://linux.oracle.com/errata/ELSA-2020-4442.html",
          "https://lists.debian.org/debian-lts-announce/2020/08/msg00037.html",
          "https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/XZARJHJJDBHI7CE5PZEBXS5HKK6HXKW2/",
          "https://security.gentoo.org/glsa/202003-16",
          "https://security.netapp.com/advisory/ntap-20190926-0003/",
          "https://security.netapp.com/advisory/ntap-20200122-0003/",
          "https://usn.ubuntu.com/4205-1/",
          "https://usn.ubuntu.com/usn/usn-4205-1",
          "https://www.mail-archive.com/sqlite-users@mailinglists.sqlite.org/msg116312.html",
          "https://www.oracle.com/security-alerts/cpuapr2020.html",
          "https://www.oracle.com/security-alerts/cpujan2020.html",
          "https://www.sqlite.org/src/info/e4598ecbdd18bd82945f6029013296690e719a62",
          "https://www.sqlite.org/src/timeline?c=98357d8c1263920b"
        ],
        "PublishedDate": "2019-09-09T17:15:00Z",
        "LastModifiedDate": "2020-08-23T01:15:00Z"
      },
      {
        "VulnerabilityID": "CVE-2019-19242",
        "PkgName": "sqlite-libs",
        "InstalledVersion": "3.26.0-r3",
        "FixedVersion": "3.28.0-r2",
        "Layer": {
          "DiffID": "sha256:fbe16fc07f0d81390525c348fbd720725dcae6498bd5e902ce5d37f2b7eed743"
        },
        "SeveritySource": "nvd",
        "PrimaryURL": "https://avd.aquasec.com/nvd/cve-2019-19242",
        "Title": "sqlite: SQL injection in sqlite3ExprCodeTarget in expr.c",
        "Description": "SQLite 3.30.1 mishandles pExpr-\u003ey.pTab, as demonstrated by the TK_COLUMN case in sqlite3ExprCodeTarget in expr.c.",
        "Severity": "MEDIUM",
        "CweIDs": [
          "CWE-476"
        ],
        "CVSS": {
          "nvd": {
            "V2Vector": "AV:N/AC:M/Au:N/C:N/I:N/A:P",
            "V3Vector": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V2Score": 4.3,
            "V3Score": 5.9
          },
          "redhat": {
            "V3Vector": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H",
            "V3Score": 5.9
          }
        },
        "References": [
          "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-19242",
          "https://github.com/sqlite/sqlite/commit/57f7ece78410a8aae86aa4625fb7556897db384c",
          "https://usn.ubuntu.com/4205-1/",
          "https://usn.ubuntu.com/usn/usn-4205-1",
          "https://www.oracle.com/security-alerts/cpuapr2020.html"
        ],
        "PublishedDate": "2019-11-27T17:15:00Z",
        "LastModifiedDate": "2020-04-15T21:15:00Z"
      }
    ]
  }
]`

	imageVuln := []*ImageVulnerabilities{}
	err := json.Unmarshal([]byte(imageVulnerabilitiesStr), &imageVuln)

	if err != nil {
		t.Errorf("Decode imageVulnerabilitiesStr failed, err: %s", err)
	} else {
		t.Log(imageVuln[0].Target)
	}
}
