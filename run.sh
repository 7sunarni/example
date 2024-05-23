(./pingtunnel -type client -l :4455 -s $SERVER -sock5 1 -ltcp :4466 -tcp_bs 4096 -lhttp :4477 -s5ftfile GeoLite2-Country.mmdb -s5filter CN)&
(./pingtunnel -type client -l :53 -s $SERVER -t 8.8.8.8:53  -nolog 1 -noprint 1 )&
