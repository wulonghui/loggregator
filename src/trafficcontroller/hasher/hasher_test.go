package hasher_test

import (
	"crypto/rand"
	"encoding/hex"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"trafficcontroller/hasher"
)

var _ = Describe("Hasher", func() {
	It("should panic when not seeded with servers", func() {
		Expect(func() {
			hasher.NewHasher([]string{})
		}).To(Panic())
	})

	Describe("LoggregatorServers", func() {

		It("should return one server", func() {
			loggregatorServer := []string{"10.10.0.16:9998"}
			h := hasher.NewHasher(loggregatorServer)
			Expect(h.LoggregatorServers()).To(Equal(loggregatorServer))
		})

		It("should return all servers", func() {
			loggregatorServers := []string{"10.10.0.16:9998", "10.10.0.17:9997"}
			h := hasher.NewHasher(loggregatorServers)
			Expect(h.LoggregatorServers()).To(Equal(loggregatorServers))
		})

	})

	Describe("GetLoggregatorServerForAppId", func() {

		It("should hashes accross one server", func() {
			loggregatorServer := []string{"10.10.0.16:9998"}
			h := hasher.NewHasher(loggregatorServer)
			ls := h.GetLoggregatorServerForAppId("app1")
			Expect(ls).To(Equal(loggregatorServer[0]))
		})

		It("should hash accross two servers", func() {
			loggregatorServer := []string{"server1", "server2"}
			h := hasher.NewHasher(loggregatorServer)
			ls := h.GetLoggregatorServerForAppId("app1")

			Expect(ls).To(Equal(loggregatorServer[1]))

			ls = h.GetLoggregatorServerForAppId("app2")
			Expect(ls).To(Equal(loggregatorServer[0]))
		})

		It("should uniformly hash traffic accross servers", func() {
			loggregatorServers := []string{"server1", "server2", "server3"}
			hitCounters := map[string]int{"server1": 0, "server2": 0, "server3": 0}

			h := hasher.NewHasher(loggregatorServers)
			target := 1000000
			for i := 0; i < target; i++ {
				ls := h.GetLoggregatorServerForAppId(GenUUID())
				hitCounters[ls] = hitCounters[ls] + 1
			}

			targetHitsPerServer := target / len(hitCounters)
			Expect(hitCounters["server1"]).To(BeNumerically("~", targetHitsPerServer, 3000))
			Expect(hitCounters["server2"]).To(BeNumerically("~", targetHitsPerServer, 3000))
			Expect(hitCounters["server3"]).To(BeNumerically("~", targetHitsPerServer, 3000))
		})

		It("should always return the same server for the given appId", func() {
			loggregatorServers := []string{"10.10.0.16:9998", "10.10.0.17:9997"}
			h := hasher.NewHasher(loggregatorServers)
			for i := 0; i < 1000; i++ {
				ls0 := h.GetLoggregatorServerForAppId("appId")
				Expect(ls0).To(Equal(loggregatorServers[0]))

				ls1 := h.GetLoggregatorServerForAppId("appId23")
				Expect(ls1).To(Equal(loggregatorServers[1]))
			}
		})
	})
})

func GenUUID() string {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		panic("No GUID generated")
	}
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid)
}
