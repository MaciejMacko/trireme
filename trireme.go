package trireme

import (
	"fmt"

	"github.com/aporeto-inc/trireme/cache"
	"github.com/aporeto-inc/trireme/enforcer"
	"github.com/aporeto-inc/trireme/policy"
	"github.com/aporeto-inc/trireme/supervisor"

	"github.com/golang/glog"
)

func addTransmitterLabel(contextID string, containerInfo *policy.PUInfo) {
	containerInfo.Policy.PolicyTags[enforcer.TransmitterLabel] = contextID
}

// trireme contains references to all the subElements of
type trireme struct {
	serverID         string
	containerTracker cache.DataStore
	supervisor       supervisor.Supervisor
	enforcer         enforcer.PolicyEnforcer
	resolver         PolicyResolver
	stopChan         chan bool
	requestChan      chan *triremeRequest
}

// NewTrireme returns a reference to the trireme object based on the parameter subelements.
func NewTrireme(serverID string, enforcer enforcer.PolicyEnforcer, supervisor supervisor.Supervisor, resolver PolicyResolver) Trireme {

	trireme := &trireme{
		serverID:         serverID,
		containerTracker: cache.NewCache(nil),
		supervisor:       supervisor,
		enforcer:         enforcer,
		resolver:         resolver,
		stopChan:         make(chan bool),
		requestChan:      make(chan *triremeRequest),
	}
	resolver.SetPolicyUpdater(trireme)

	return trireme
}

// Start starts trireme individual components.
func (t *trireme) Start() error {

	if err := t.supervisor.Start(); err != nil {
		return fmt.Errorf("Error starting Controller: %s", err)
	}

	if err := t.enforcer.Start(); err != nil {
		return fmt.Errorf("Error starting enforcer: %s", err)
	}

	// Starting main trireme routine
	go t.triremeWorker()

	return nil
}

// Stop stops trireme individual components
func (t *trireme) Stop() error {

	// send the stop signal for the trireme worker routine.
	t.stopChan <- true

	if err := t.supervisor.Stop(); err != nil {
		return fmt.Errorf("Error stopping Controller: %s", err)
	}

	if err := t.enforcer.Stop(); err != nil {
		return fmt.Errorf("Error stopping enforcer: %s", err)
	}

	return nil
}

// HandleCreate is acting on a create monitoring event.
func (t *trireme) HandleCreate(contextID string, runtimeInfo *policy.PURuntime) <-chan error {

	c := make(chan error)

	req := &triremeRequest{
		contextID:   contextID,
		reqType:     requestCreate,
		runtimeInfo: runtimeInfo,
		returnChan:  c,
	}

	t.requestChan <- req

	return c
}

// HandleDelete is acting on a delete monitoring object
func (t *trireme) HandleDelete(contextID string) <-chan error {

	c := make(chan error)

	req := &triremeRequest{
		contextID:  contextID,
		reqType:    requestDelete,
		returnChan: c,
	}

	t.requestChan <- req

	return c
}

func (t *trireme) UpdatePolicy(contextID string, newPolicy *policy.PUPolicy) <-chan error {

	c := make(chan error)

	req := &triremeRequest{
		contextID:  contextID,
		reqType:    policyUpdate,
		policyInfo: newPolicy,
		returnChan: c,
	}

	t.requestChan <- req

	return c
}

func (t *trireme) PURuntime(contextID string) (policy.RuntimeReader, error) {

	container, err := t.containerTracker.Get(contextID)

	if err != nil {
		return nil, err
	}

	return container.(*policy.PURuntime), nil
}

func (t *trireme) doHandleCreate(contextID string, runtimeInfo *policy.PURuntime) error {

	// Cache all the container runtime information
	if err := t.containerTracker.AddOrUpdate(contextID, runtimeInfo); err != nil {
		return err
	}

	policyInfo, err := t.resolver.GetPolicy(contextID, runtimeInfo)
	if err != nil {
		glog.V(2).Infoln("Policy Error for this context: %s . Container killed. %s", contextID, err)
		return fmt.Errorf("Policy Error for this context: %s . Container killed. %s", contextID, err)
	}

	containerInfo := policy.PUInfoFromPolicyAndRuntime(contextID, policyInfo, runtimeInfo)

	addTransmitterLabel(contextID, containerInfo)

	err = t.supervisor.AddPU(contextID, containerInfo)
	if err != nil {
		t.resolver.DeletePU(contextID)
		return fmt.Errorf("Not able to setup supervisor: %s", err)
	}

	err = t.enforcer.AddPU(contextID, containerInfo)
	if err != nil {
		t.resolver.DeletePU(contextID)
		t.supervisor.DeletePU(contextID)
		return fmt.Errorf("Not able to setup enforcer: %s", err)
	}
	glog.V(2).Infoln("Finished PUHandleCreate: %s .", contextID)
	return nil
}

func (t *trireme) doHandleDelete(contextID string) error {
	t.resolver.DeletePU(contextID)
	t.supervisor.DeletePU(contextID)

	runtimeInfo, err := t.PURuntime(contextID)
	t.containerTracker.Remove(contextID)
	if err != nil {
		return fmt.Errorf("Error getting Runtime out of cache for ContextID %s : %s", contextID, err)
	}
	ip, ok := runtimeInfo.DefaultIPAddress()
	if !ok {
		return fmt.Errorf("default IPAddress not found for %s", contextID)
	}
	t.enforcer.DeletePU(ip)
	glog.V(5).Infof("Finished HandleDelete. %s", contextID)
	return nil
}

func (t *trireme) doUpdatePolicy(contextID string, newPolicy *policy.PUPolicy) error {
	runtimeInfo, err := t.PURuntime(contextID)
	if err != nil {
		return err
	}
	containerInfo := policy.PUInfoFromPolicyAndRuntime(contextID, newPolicy, runtimeInfo.(*policy.PURuntime))

	addTransmitterLabel(contextID, containerInfo)

	err = t.supervisor.UpdatePU(contextID, containerInfo)
	if err != nil {
		return err
	}

	err = t.enforcer.UpdatePU(containerInfo.Runtime.IPAddresses()["bridge"], containerInfo)
	if err != nil {
		t.supervisor.DeletePU(contextID)
		return err
	}
	glog.V(5).Infof("Finished UpdatePolicy. %s", contextID)
	return nil
}

func (t *trireme) handleRequest(request *triremeRequest) error {
	switch request.reqType {
	case requestCreate:
		return t.doHandleCreate(request.contextID, request.runtimeInfo)
	case requestDelete:
		return t.doHandleDelete(request.contextID)
	case policyUpdate:
		return t.doUpdatePolicy(request.contextID, request.policyInfo)
	default:
		return fmt.Errorf("trireme Request format not recognized: %d", request.reqType)
	}
}

func (t *trireme) triremeWorker() {
	for {
		select {
		case <-t.stopChan:
			glog.V(2).Infof("Stopping trireme worker.")
			return
		case req := <-t.requestChan:
			glog.V(5).Infof("Handling trireme Request Type %d ", req.reqType)
			req.returnChan <- t.handleRequest(req)
		}
	}
}
