#  Gereksinimler :
- Oyun sunucusu oyuncuların kayıt olmak için kullanabilecekleri bir post request
sağlayacaktır. Bu post request bodysinde bir nickname alır ve random bir kullanıcı idsi
döner. Dönülen id ileride yapılacak websocket requestlerinde kullanılır ve sisteme
kaydolmamış idler notRegistered tipinde bir hata alır.
- Basit bir matchmaking algoritması kullanarak websocket üzerinden kabul ettiği oyuncu
isteklerini kabul edip, 30 saniyede bir içeride biriken kullanıcıları 3 lü gruplar halinde
eşleştirerek birer oda oluşturacak ve odanın idsi daha sonraki komutlarda kullanılmak
üzere oyunculara iletilecek
- Her oda oyuncuların bilmediği bir sayı içerecek, ve oyuncular bu sayıyı tahmin etmek için
bir websocket komutu gönderebilecekler
- Bir odadaki tüm oyuncular tahminlerini gönderdiklerinde en yakın tahminde bulunan
oyuncu oyunu kazanacak ve kazanan oyuncu, tahmini, gizli sayı ve tahminlere göre
oyuncuların sırası oyunculara yine websocket üzerinden iletilecek
- Eğer bir oyuncu 20 saniye boyunca bir tahmin iletmezse oyuncu otomatik olarak
kaybedecek
- Oyun sunucusu içerideki odalar ve kullanıcılar ile ilgili istatistiklerin sorgulanabileceği bir
get isteği sağlayacak. Bu get isteği kayıtlı kullanıcı sayısını ve içerideki aktif odalar ile
ilgili bilgileri içerecek
- Hem websocket hem de web istekleri için kullanılacak json protokolu aşağıda verilmiştir


**JSON Protokolü:**

**/register:** Bu bir post isteğidir. Kullanıcı bu istek bodysinde içinde nickname olan
bir json gönderir ve kullanıcı idsinin olduğu bir json cevabı alır. Bu request sonucunda sunucu
memoryde kupa, nickname ve idden oluşan bir oyuncu kaydı oluşturup saklar.
Request: {“nickname”:string}
Response: {“id”: uuid}

**/stats:** Bu bir get isteğidir. Bu requestin cevabında oyuna kaydolmuş kullanıcı sayısı,
aktif odaların gizli sayılar ve idleri bulunur.
Response: {“registeredPlayers”:int, “activeRooms”: [
{“id”: int, “secret”: int}, ....
]}

**join**: Bu bir websocket komutudur. Bu komut kullanıcının idsini alır ve eğer kullanıcı
kayıtlı bir kullanıcı ise sunucuda bekleyen istekler listesine ekler ve kullanıcıya eklendiğine dair
bir cevap döner. Eğer kullanıcı kayıtlı kullanıcılar arasında değilse notRegistered hatası döner.
30 saniyede bir bekleyen isteklerdeki kullanıcılar kupalarına göre eşleştirilerek bir oda
oluşturulur ve istek gönderen kullanıcılara eşleştirildikleri oda idsi gönderilir.
Command json:{“cmd”:”join”, “id”:uuid}
Reply json: {“cmd”: “join”, “reply”:”waiting”}
Error json: {“cmd”:”join”, “error”: “notRegistered”}
joinedRoom: Bu bir websocket eventidir. Join komutunu gönderen kullanıcıya eşleşme
yapılıp bir odaya eklendiği zaman gönderilir
Event json: {“event”:”joinedRoom”, “room”: int}

**guess:** Bu bir websocket komutudur. Oyuncu odaya katıldıktan sonra idsi, odası ve
tahminini içeren guess komutunu gönderebilir. Oyuncu kayıtlı değilse notRegistered, oyuncu
bu odada değilse notInRoom hatası döner. Odaya katıldıktan sonra bu komutu göndermek için
20 saniyesi olacaktır bu süre içinde tahmin göndermezse oyuncu timeout olur ve oyun sonlanır.
Command json: {“cmd“: “guess”, “user”:uuid, “room”:int, “data”:int}
Reply json: {“cmd”:”guess”, “reply”:”guessReceived”}
Error json: {“cmd”:”join”, “error”: “notRegistered”}
Error json: {“cmd”:”join”, “error”: “notInRoom”}

**gameOver:** Bu bir websocket eventidir. Bir odada tüm tahminler alındıktan veya 20
saniye geçtikten sonra, odanın gizli sayısı ve alınan tahminlerin odadaki gizli sayıya yakınlığına
göre oyuncuların sıralamasını içerir. 1. oyuncu 30 kupa kazanır, 2. oyuncu 10 kupa ve 3. oyuncu
kupa kazanamaz. Tahmin göndermeyen kullanıcılar da kupa kazanmazlar. Eğer bir oyuncu
tahmin gönderemediyse sırası -1 olarak gönderilir
Event json: {“event”:”gameOver”, “secret”:int, “rankings”: [
{“rank”:1, “player”:uuid, “guess”:int, “deltaTrophy”:30},
{“rank”:2, “player”:uuid},”deltaTrophy”:20},
{“rank”:3, “player”:uuid},”deltaTrophy”:0}]}
