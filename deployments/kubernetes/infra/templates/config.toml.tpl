port = ${port}
env = "${environment}"

[db]
max-open-conns = "${max_open_conns}"
max-idle-conns = "${max_idle_conns}"
max-idle-time = "${max_idle_time}"

[limiter]
rps = ${rps}
burst = ${burst}
enabled = ${enabled}
