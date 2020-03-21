/*global require,setInterval,console */
const opcua = require("node-opcua");

// Let's create an instance of OPCUAServer
const server = new opcua.OPCUAServer({
    port: 4334, // the port of the listening socket of the server
    resourcePath: "/", // this path will be added to the endpoint resource name
     buildInfo : {
        productName: "MySampleServer1",
        buildNumber: "7658",
        buildDate: new Date(2014,5,2)
    }
});

function post_initialize() {
    console.log("initialized");
    function construct_my_address_space(server) {
    
        const addressSpace = server.engine.addressSpace;
        const namespace = addressSpace.getOwnNamespace();
    
        // declare a new object
        const device = namespace.addObject({
            organizedBy: addressSpace.rootFolder.objects,
            browseName: "MyDevice"
        });
    
        // add some variables
        // add a variable named MyVariable1 to the newly created folder "MyDevice"
        let PRESSURE = 1;
        let HUMIDITY = 2;
        let TEMPHOME = 3;
        let TEMPWORK = 4;
        let LEVELPH = 5;
        let MASS = 6;
        let WATER = 7;
        let LEVELCO2 = 8;
        // emulate variable1 changing every 500 ms
        //setInterval(function(){  variable1+=1; }, 500);
        
        namespace.addVariable({
            nodeId: "s=pressure",
            componentOf: device,
            browseName: "pressure",
            dataType: "Double",
            value: {
                get: function () {
                    return new opcua.Variant({dataType: opcua.DataType.Double, value: PRESSURE });
                }
            }
        });
        namespace.addVariable({
            nodeId: "s=humidity",
            componentOf: device,
            browseName: "humidity",
            dataType: "Double",
            value: {
                get: function () {
                    return new opcua.Variant({dataType: opcua.DataType.Double, value: HUMIDITY });
                }
            }
        });
        namespace.addVariable({
            nodeId: "s=temphome",
            componentOf: device,
            browseName: "temphome",
            dataType: "Double",
            value: {
                get: function () {
                    return new opcua.Variant({dataType: opcua.DataType.Double, value: TEMPHOME });
                }
            }
        });
        namespace.addVariable({
            nodeId: "s=tempwork",
            componentOf: device,
            browseName: "tempwork",
            dataType: "Double",
            value: {
                get: function () {
                    return new opcua.Variant({dataType: opcua.DataType.Double, value: TEMPWORK });
                }
            }
        });
        namespace.addVariable({
            nodeId: "s=levelph",
            componentOf: device,
            browseName: "levelph",
            dataType: "Double",
            value: {
                get: function () {
                    return new opcua.Variant({dataType: opcua.DataType.Double, value: LEVELPH });
                }
            }
        });
        namespace.addVariable({
            nodeId: "s=levelco2",
            componentOf: device,
            browseName: "levelco2",
            dataType: "Double",
            value: {
                get: function () {
                    return new opcua.Variant({dataType: opcua.DataType.Double, value: LEVELCO2 });
                }
            }
        });
        namespace.addVariable({
            nodeId: "s=mass",
            componentOf: device,
            browseName: "mass",
            dataType: "Double",
            value: {
                get: function () {
                    return new opcua.Variant({dataType: opcua.DataType.Double, value: MASS });
                }
            }
        });
        namespace.addVariable({
            nodeId: "s=water",
            componentOf: device,
            browseName: "water",
            dataType: "Double",
            value: {
                get: function () {
                    return new opcua.Variant({dataType: opcua.DataType.Double, value: WATER });
                }
            }
        });
        const os = require("os");
        /**
         * returns the percentage of free memory on the running machine
         * @return {double}
         */
        function available_memory() {
            // var value = process.memoryUsage().heapUsed / 1000000;
            const percentageMemUsed = os.freemem() / os.totalmem() * 100.0;
            return percentageMemUsed;
        }

    }
    construct_my_address_space(server);
    server.start(function() {
        console.log("Server is now listening ... ( press CTRL+C to stop)");
        console.log("port ", server.endpoints[0].port);
        const endpointUrl = server.endpoints[0].endpointDescriptions()[0].endpointUrl;
        console.log(" the primary server endpoint url is ", endpointUrl );
    });
}
server.initialize(post_initialize);