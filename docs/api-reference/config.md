<p>Packages:</p>
<ul>
<li>
<a href="#config.trust-configurator.gardener.cloud%2fv1alpha1">config.trust-configurator.gardener.cloud/v1alpha1</a>
</li>
</ul>

<h2 id="config.trust-configurator.gardener.cloud/v1alpha1">config.trust-configurator.gardener.cloud/v1alpha1</h2>
<p>

</p>

<h3 id="controllerconfiguration">ControllerConfiguration
</h3>


<p>
(<em>Appears on:</em><a href="#gardenshoottrustconfiguratorconfiguration">GardenShootTrustConfiguratorConfiguration</a>)
</p>

<p>
ControllerConfiguration defines the configuration of the controllers.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>shoot</code></br>
<em>
<a href="#shootcontrollerconfig">ShootControllerConfig</a>
</em>
</td>
<td>
<p>Shoot is the configuration for the shoot controller.</p>
</td>
</tr>
<tr>
<td>
<code>garbageCollector</code></br>
<em>
<a href="#garbagecollectorcontrollerconfig">GarbageCollectorControllerConfig</a>
</em>
</td>
<td>
<p>GarbageCollector is the configuration for the garbage-collector controller.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="garbagecollectorcontrollerconfig">GarbageCollectorControllerConfig
</h3>


<p>
(<em>Appears on:</em><a href="#controllerconfiguration">ControllerConfiguration</a>)
</p>

<p>
GarbageCollectorControllerConfig is the configuration for the garbage-collector controller.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>syncPeriod</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta">Duration</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SyncPeriod is the duration how often the controller performs its reconciliation.</p>
</td>
</tr>
<tr>
<td>
<code>minimumObjectLifetime</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta">Duration</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>MinimumObjectLifetime is the minimum age an object must have before it is considered for garbage collection.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="gardenshoottrustconfiguratorconfiguration">GardenShootTrustConfiguratorConfiguration
</h3>


<p>
GardenShootTrustConfiguratorConfiguration defines the configuration for the Gardener garden-shoot-trust-configurator.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>kind</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds</p>
</td>
</tr>
<tr>
<td>
<code>apiVersion</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources</p>
</td>
</tr>
<tr>
<td>
<code>leaderElection</code></br>
<em>
<a href="#leaderelectionconfiguration">LeaderElectionConfiguration</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>LeaderElection defines the configuration of leader election client.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code></br>
<em>
string
</em>
</td>
<td>
<p>LogLevel is the level/severity for the logs. Must be one of [info,debug,error].</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code></br>
<em>
string
</em>
</td>
<td>
<p>LogFormat is the output format for the logs. Must be one of [text,json].</p>
</td>
</tr>
<tr>
<td>
<code>controllers</code></br>
<em>
<a href="#controllerconfiguration">ControllerConfiguration</a>
</em>
</td>
<td>
<p>Controllers defines the configuration of the controllers.</p>
</td>
</tr>
<tr>
<td>
<code>server</code></br>
<em>
<a href="#serverconfiguration">ServerConfiguration</a>
</em>
</td>
<td>
<p>Server defines the configuration of the HTTP server.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="httpsserver">HTTPSServer
</h3>


<p>
(<em>Appears on:</em><a href="#serverconfiguration">ServerConfiguration</a>)
</p>

<p>
HTTPSServer is the configuration for the HTTPSServer server.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>port</code></br>
<em>
integer
</em>
</td>
<td>
<p>Port is the port on which to serve requests.</p>
</td>
</tr>
<tr>
<td>
<code>bindAddress</code></br>
<em>
string
</em>
</td>
<td>
<p>BindAddress is the IP address on which to listen for the specified port.</p>
</td>
</tr>
<tr>
<td>
<code>tls</code></br>
<em>
<a href="#tls">TLS</a>
</em>
</td>
<td>
<p>TLS contains information about the TLS configuration for a HTTPS server.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="oidcconfig">OIDCConfig
</h3>


<p>
(<em>Appears on:</em><a href="#shootcontrollerconfig">ShootControllerConfig</a>)
</p>

<p>
OIDCConfig is the configuration for the OIDC resources created for trusted shoots.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>audiences</code></br>
<em>
string array
</em>
</td>
<td>
<em>(Optional)</em>
<p>Audiences is the list of audience identifiers used in the OIDC resources for trusted shoots.<br />Defaults to ["garden"].</p>
</td>
</tr>
<tr>
<td>
<code>maxTokenExpiration</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta">Duration</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxTokenExpiration sets a limit to the maximum validity duration of a token.<br />Tokens issued with validity greater than this value will not be verified.<br />Must be between 5 minutes and 24 hours. Defaults to 2 hours.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="server">Server
</h3>


<p>
(<em>Appears on:</em><a href="#httpsserver">HTTPSServer</a>, <a href="#serverconfiguration">ServerConfiguration</a>)
</p>

<p>
Server contains information for HTTP(S) server configuration.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>port</code></br>
<em>
integer
</em>
</td>
<td>
<p>Port is the port on which to serve requests.</p>
</td>
</tr>
<tr>
<td>
<code>bindAddress</code></br>
<em>
string
</em>
</td>
<td>
<p>BindAddress is the IP address on which to listen for the specified port.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="serverconfiguration">ServerConfiguration
</h3>


<p>
(<em>Appears on:</em><a href="#gardenshoottrustconfiguratorconfiguration">GardenShootTrustConfiguratorConfiguration</a>)
</p>

<p>
ServerConfiguration contains details for the HTTP(S) servers.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>webhooks</code></br>
<em>
<a href="#httpsserver">HTTPSServer</a>
</em>
</td>
<td>
<p>Webhooks is the configuration for the HTTPS webhook server.</p>
</td>
</tr>
<tr>
<td>
<code>healthProbes</code></br>
<em>
<a href="#server">Server</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>HealthProbes is the configuration for serving the healthz and readyz endpoints.</p>
</td>
</tr>
<tr>
<td>
<code>metrics</code></br>
<em>
<a href="#server">Server</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Metrics is the configuration for serving the metrics endpoint.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="shootcontrollerconfig">ShootControllerConfig
</h3>


<p>
(<em>Appears on:</em><a href="#controllerconfiguration">ControllerConfiguration</a>)
</p>

<p>
ShootControllerConfig is the configuration for the shoot controller.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>syncPeriod</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta">Duration</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SyncPeriod is the duration how often the controller performs its reconciliation.</p>
</td>
</tr>
<tr>
<td>
<code>oidcConfig</code></br>
<em>
<a href="#oidcconfig">OIDCConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>OIDCConfig is the configuration for the OIDC resources which are created for trusted shoots.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="tls">TLS
</h3>


<p>
(<em>Appears on:</em><a href="#httpsserver">HTTPSServer</a>)
</p>

<p>
TLS contains information about the TLS configuration for a HTTPS server.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>serverCertDir</code></br>
<em>
string
</em>
</td>
<td>
<p>ServerCertDir is the path to a directory containing the server's TLS certificate and key (the files must be<br />named tls.crt and tls.key respectively).</p>
</td>
</tr>

</tbody>
</table>


