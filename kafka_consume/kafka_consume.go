package kafka_consume

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Hymiside/evraz-api/pkg/service"
	"github.com/Shopify/sarama"
)

func InitConsume(s *service.Service) {
	brokers := "rc1a-b5e65f36lm3an1d5.mdb.yandexcloud.net:9091"
	splitBrokers := strings.Split(brokers, ",")
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Version = sarama.V0_10_0_0
	conf.Consumer.Return.Errors = true
	conf.ClientID = "sasl_scram_client"
	conf.Metadata.Full = true
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = "9433_reader"
	conf.Net.SASL.Password = "eUIpgWu0PWTJaTrjhjQD3.hoyhntiK"
	conf.Net.SASL.Handshake = true
	conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	conf.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA512)

	certs := x509.NewCertPool()
	pemPath := "/usr/local/share/ca-certificates/Yandex/YandexCA.crt"
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		fmt.Println("Couldn't load cert: ", err.Error())
		// Handle the error
	}
	certs.AppendCertsFromPEM(pemData)

	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            certs,
	}

	master, err := sarama.NewConsumer(splitBrokers, conf)
	if err != nil {
		fmt.Println("Coulnd't create consumer: ", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err = master.Close(); err != nil {
			panic(err)
		}
	}()

	topic := "zsmk-9433-dev-01"

	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	var m map[string]interface{}

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err = <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				if err = json.Unmarshal(msg.Value, &m); err != nil {
					panic(err)
				}
				//fmt.Println(m)
				s.Evraz.Consume(m)
			}
		}
	}()
	<-doneCh
}
