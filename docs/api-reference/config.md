<p>Packages:</p>
<ul>
<li>
<a href="#config.trust-configurator.gardener.cloud%2fv1alpha1">config.trust-configurator.gardener.cloud/v1alpha1</a>
</li>
</ul>
<h2 id="config.trust-configurator.gardener.cloud/v1alpha1">config.trust-configurator.gardener.cloud/v1alpha1</h2>
<p>
<p>Package v1alpha1 contains the shoot trust confogurator configuration.</p>
</p>
Resource Types:
<ul></ul>
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.ControllerConfiguration">ControllerConfiguration
</h3>
<p>
(<em>Appears on:</em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.GardenShootTrustConfiguratorConfiguration">GardenShootTrustConfiguratorConfiguration</a>)
</p>
<p>
<p>ControllerConfiguration defines the configuration of the controllers.</p>
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
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.ShootControllerConfig">
ShootControllerConfig
</a>
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
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.GarbageCollectorControllerConfig">
GarbageCollectorControllerConfig
</a>
</em>
</td>
<td>
<p>GarbageCollector is the configuration for the garbage-collector controller.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.GarbageCollectorControllerConfig">GarbageCollectorControllerConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.ControllerConfiguration">ControllerConfiguration</a>)
</p>
<p>
<p>GarbageCollectorControllerConfig is the configuration for the garbage-collector controller.</p>
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#duration-v1-meta">
Kubernetes meta/v1.Duration
</a>
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#duration-v1-meta">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>MinimumObjectLifetime is the minimum age an object must have before it is considered for garbage collection.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.GardenShootTrustConfiguratorConfiguration">GardenShootTrustConfiguratorConfiguration
</h3>
<p>
<p>GardenShootTrustConfiguratorConfiguration defines the configuration for the Gardener garden-shoot-trust-configurator.</p>
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
<code>leaderElection</code></br>
<em>
k8s.io/component-base/config/v1alpha1.LeaderElectionConfiguration
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
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.ControllerConfiguration">
ControllerConfiguration
</a>
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
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.ServerConfiguration">
ServerConfiguration
</a>
</em>
</td>
<td>
<p>Server defines the configuration of the HTTP server.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.HTTPSServer">HTTPSServer
</h3>
<p>
(<em>Appears on:</em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.ServerConfiguration">ServerConfiguration</a>)
</p>
<p>
<p>HTTPSServer is the configuration for the HTTPSServer server.</p>
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
<code>Server</code></br>
<em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.Server">
Server
</a>
</em>
</td>
<td>
<p>
(Members of <code>Server</code> are embedded into this type.)
</p>
<p>Server is the configuration for the bind address and the port.</p>
</td>
</tr>
<tr>
<td>
<code>tls</code></br>
<em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.TLS">
TLS
</a>
</em>
</td>
<td>
<p>TLS contains information about the TLS configuration for a HTTPS server.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.OIDCConfig">OIDCConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.ShootControllerConfig">ShootControllerConfig</a>)
</p>
<p>
<p>OIDCConfig is the configuration for the OIDC resources created for trusted shoots.</p>
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
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Audiences is the list of audience identifiers used in the OIDC resources for trusted shoots.
Defaults to [&ldquo;garden&rdquo;].</p>
</td>
</tr>
<tr>
<td>
<code>maxTokenExpiration</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#duration-v1-meta">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxTokenExpiration sets a limit to the maximum validity duration of a token.
Tokens issued with validity greater than this value will not be verified.
Must be between 5 minutes and 24 hours. Defaults to 2 hours.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.Server">Server
</h3>
<p>
(<em>Appears on:</em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.HTTPSServer">HTTPSServer</a>, 
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.ServerConfiguration">ServerConfiguration</a>)
</p>
<p>
<p>Server contains information for HTTP(S) server configuration.</p>
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
int
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
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.ServerConfiguration">ServerConfiguration
</h3>
<p>
(<em>Appears on:</em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.GardenShootTrustConfiguratorConfiguration">GardenShootTrustConfiguratorConfiguration</a>)
</p>
<p>
<p>ServerConfiguration contains details for the HTTP(S) servers.</p>
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
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.HTTPSServer">
HTTPSServer
</a>
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
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.Server">
Server
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>HealthProbes is the configuration for serving the healthz and readyz endpoints.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.ShootControllerConfig">ShootControllerConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.ControllerConfiguration">ControllerConfiguration</a>)
</p>
<p>
<p>ShootControllerConfig is the configuration for the shoot controller.</p>
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
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#duration-v1-meta">
Kubernetes meta/v1.Duration
</a>
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
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.OIDCConfig">
OIDCConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>OIDCConfig is the configuration for the OIDC resources which are created for trusted shoots.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="config.trust-configurator.gardener.cloud/v1alpha1.TLS">TLS
</h3>
<p>
(<em>Appears on:</em>
<a href="#config.trust-configurator.gardener.cloud/v1alpha1.HTTPSServer">HTTPSServer</a>)
</p>
<p>
<p>TLS contains information about the TLS configuration for a HTTPS server.</p>
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
<p>ServerCertDir is the path to a directory containing the server&rsquo;s TLS certificate and key (the files must be
named tls.crt and tls.key respectively).</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <a href="https://github.com/ahmetb/gen-crd-api-reference-docs">gen-crd-api-reference-docs</a>
</em></p>
