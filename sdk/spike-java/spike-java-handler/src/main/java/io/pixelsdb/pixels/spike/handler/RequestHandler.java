package io.pixelsdb.pixels.spike.handler;

public interface RequestHandler {
    SpikeWorker.CallWorkerFunctionResp execute(SpikeWorker.CallWorkerFunctionReq request);
}
