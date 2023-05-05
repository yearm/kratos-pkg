package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Prom ...
type Prom struct {
	namespace string
	histogram *prometheus.HistogramVec
	counter   *prometheus.CounterVec
	summary   *prometheus.SummaryVec
	gauge     *prometheus.GaugeVec
}

// NewProm ...
func NewProm(namespace string) *Prom {
	return &Prom{namespace: namespace}
}

// RegisterHistogram ...
func (p *Prom) RegisterHistogram(name string, labels []string) *Prom {
	if p == nil || p.histogram != nil {
		return p
	}
	p.histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: p.namespace,
		Name:      name,
	}, labels)
	prometheus.MustRegister(p.histogram)
	return p
}

// HistogramObserve ...
func (p *Prom) HistogramObserve(v float64, labels ...string) {
	if p.histogram != nil {
		p.histogram.WithLabelValues(labels...).Observe(v)
	}
}

// RegisterCounter ...
func (p *Prom) RegisterCounter(name string, labels []string) *Prom {
	if p == nil || p.counter != nil {
		return p
	}
	p.counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: p.namespace,
		Name:      name,
	}, labels)
	prometheus.MustRegister(p.counter)
	return p
}

// CounterIncr ...
func (p *Prom) CounterIncr(labels ...string) {
	if p.counter != nil {
		p.counter.WithLabelValues(labels...).Inc()
	}
}

// CounterAdd ...
func (p *Prom) CounterAdd(v float64, labels ...string) {
	if p.counter != nil {
		p.counter.WithLabelValues(labels...).Add(v)
	}
}

// RegisterSummary ...
func (p *Prom) RegisterSummary(name string, labels []string) *Prom {
	if p == nil || p.summary != nil {
		return p
	}
	p.summary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: p.namespace,
		Name:      name,
	}, labels)
	prometheus.MustRegister(p.summary)
	return p
}

// SummaryObserve ...
func (p *Prom) SummaryObserve(v float64, labels ...string) {
	if p.summary != nil {
		p.summary.WithLabelValues(labels...).Observe(v)
	}
}

// RegisterGauge ...
func (p *Prom) RegisterGauge(name string, labels []string) *Prom {
	if p == nil || p.gauge != nil {
		return p
	}
	p.gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: p.namespace,
		Name:      name,
	}, labels)
	prometheus.MustRegister(p.gauge)
	return p
}

// GaugeAdd ...
func (p *Prom) GaugeAdd(v float64, labels ...string) {
	if p.gauge != nil {
		p.gauge.WithLabelValues(labels...).Add(v)
	}
}

// GaugeInc ...
func (p *Prom) GaugeInc(labels ...string) {
	if p.gauge != nil {
		p.gauge.WithLabelValues(labels...).Inc()
	}
}

// GaugeSet ...
func (p *Prom) GaugeSet(v float64, labels ...string) {
	if p.gauge != nil {
		p.gauge.WithLabelValues(labels...).Set(v)
	}
}
