topic-cmd := kafka-topics --zookeeper localhost:2181
groups-cmd := kafka-consumer-groups --bootstrap-server localhost:9092

.PHONY: create-topics delete-topics describe-topics reset-offsets describe-groups


reset: delete-topics create-topics reset-offsets


create-topics:
	# Tracker topic is automatically created with default settings on produce
	$(topic-cmd) --create --topic worker-t --partitions 3 --replication-factor 1

delete-topics:
	$(topic-cmd) --delete --topic tracker-t
	$(topic-cmd) --delete --topic worker-t

describe-topics:
	$(topic-cmd) --describe



reset-offsets:
	$(groups-cmd) --reset-offsets --to-offset 0 --topic tracker-t --group tracker-g
	$(groups-cmd) --reset-offsets --to-offset 0 --topic worker-t  --group worker-g

describe-groups:
	$(groups-cmd) --describe --group tracker-g
	$(groups-cmd) --describe --group worker-g
