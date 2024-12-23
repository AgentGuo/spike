package io.pixelsdb.pixels.spike.handler;

import io.grpc.stub.StreamObserver;

public class FunctionServiceImpl extends SpikeWorkerServiceGrpc.SpikeWorkerServiceImplBase {

    private final RequestHandler handler;

    public FunctionServiceImpl(RequestHandler handler) {
        this.handler = handler;
    }

    @Override
    public void callWorkerFunction(SpikeWorker.CallWorkerFunctionReq request, StreamObserver<SpikeWorker.CallWorkerFunctionResp> responseObserver) {
        try {
            SpikeWorker.CallWorkerFunctionResp resp = handler.execute(request);
            responseObserver.onNext(resp);
            responseObserver.onCompleted();
        } catch (Exception e) {
            responseObserver.onError(e);
        }
    }
}