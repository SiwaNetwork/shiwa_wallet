<template>
  <div id="app" class="container">
    <div class="row">
      <div class="col-md-4 d-flex">
        <stats-card class="stats-card">
          <div class="icon-big text-center" :class="`icon-success`" slot="header">
            <i class="fas fa-wallet"></i>
          </div>
          <div class="numbers text-center text-overflow" slot="content">
            <p>SHI</p>
            {{tokenAmount | asCurrency('SHI')}}
          </div>
          <div class="stats" slot="footer">
            <i class="ti-timer"></i>
            Updates in {{counters.token}} seconds
          </div>
        </stats-card>
      </div>
      <div class="col-md-4 d-flex">
        <stats-card class="stats-card">
          <div class="icon-big text-center" :class="`icon-danger`" slot="header">
            <i class="fas fa-search-dollar"></i>
          </div>
          <div class="numbers text-center text-overflow" slot="content">
            <p>{{currency}}</p>
            {{totalValue | asCurrency(currency)}}
          </div>
          <div class="stats" slot="footer">
            <i class="ti-timer"></i>
            Updates in {{counters.value}} seconds
          </div>
        </stats-card>
      </div>
      <div class="col-md-4 d-flex">
        <stats-card class="stats-card">
          <div class="icon-big text-center" :class="`icon-info`" slot="header">
            <i class="fas fa-cube"></i>
          </div>
          <div class="numbers text-center text-overflow" slot="content">
            <p>Blocks</p>
            {{stat.blocks}}
            <!-- {{this.$store.state.OS.windows}}
            {{this.$store.state.OS.macOS}}
            {{this.$store.state.OS.linux}} -->
          </div>
          <div class="stats" slot="footer">
            <i class="ti-reload"></i>
            Updates in {{counters.block}} seconds
          </div>
        </stats-card>
      </div>
    </div>
    <div class="row">
      <div class="col">
        <card title="Wallet Address" sub-title>
          <div class="wallet-address">
            <table style="table-layout:fixed;" class="table-noheader">
              <tr>
                <td
                  class="text-overflow"
                  style="word-wrap:break-word; padding-top: 20px; padding-left: 15px; width: 100%;"
                >
                  <span style="width: 100%;" class="text-overflow">{{address}}</span>
                  <input type="hidden" id="testing-code" :value="address" />
                </td>
                <td style="padding-top: 10px; width: 9%;">
                  <p-button
                    type="info"
                    style="margin-bottom: 5px;"
                    icon
                    @click.native="copyTestingCode"
                  >
                    <i class="fa fa-copy"></i>
                  </p-button>
                </td>
              </tr>
            </table>
          </div>
        </card>
      </div>
    </div>
    <div class="row">
      <div v-if="toggle.nodesOnline" class="col-md-6 col-12 d-flex">
        <chart-card
          title="Nodes Online"
          sub-title="Since last 24 hours"
          :chart-data="chart.nodesOnline"
          chart-type="Pie"
        >
          <div slot="legend">
            <i class="fa fa-circle text-info"></i> Foundation
            <i class="fa fa-circle text-success"></i> Medium
            <i class="fa fa-circle text-secondary"></i> Light
          </div>
          <span slot="footer">
            <i class="ti-timer"></i>
            Updates in {{counters.chart}} seconds
          </span>
        </chart-card>
      </div>

      <div v-if="toggle.transactions" class="col-md-6 col-12 d-flex">
        <chart-card
          title="Transactions"
          sub-title="Sent vs. received over the last year"
          :chart-data="chart.transactions"
          :chart-options="transactionChart.options"
          chart-type="Line"
        >
          <div slot="legend">
            <i class="fa fa-circle text-info"></i> TX
            <i class="fa fa-circle text-success"></i> RX
          </div>
          <span style="padding-top: 0.625em;" slot="footer">
            <i class="ti-timer"></i>
            Updates in {{counters.chart}} seconds
          </span>
        </chart-card>
      </div>

      <div v-if="toggle.throughput" class="col-md-6 col-12 d-flex">
        <chart-card
          title="Network Throughput (tps)"
          sub-title="24 Hours performance"
          :chart-data="chart.throughput"
          :chart-options="this.throughputChart.options"
          chart-type="Line"
        >
          <span slot="footer">
            <i class="ti-timer"></i>
            Updates in {{counters.chart}} seconds
          </span>
          <div slot="legend">
            <i class="fa fa-circle text-info"></i> $SHI Tokens
            <i class="fa fa-circle text-success"></i> Data
          </div>
        </chart-card>
      </div>
    </div>
  </div>
</template>

<script>
import {mapState} from 'vuex'
import { StatsCard, ChartCard} from "@/components/index";
import Chartist from 'chartist';
import WalletCopiedNotification from "./Notifications/WalletCopied";
import WalletCopiedFailedNotification from "./Notifications/WalletCopiedFailed";

export default {
  components: {
    StatsCard,
    ChartCard
  },
  methods: {
    copyTestingCode() {
      let testingCodeToCopy = document.querySelector("#testing-code");
      testingCodeToCopy.setAttribute("type", "text");
      testingCodeToCopy.select();

      try {
        var successful = document.execCommand("copy");
        successful ? "successful" : "unsuccessful";
        this.addNotification("top", "right", 2, WalletCopiedNotification);
      } catch (err) {
        this.addNotification("top", "right", 4, WalletCopiedFailedNotification);
        alert("Oops, unable to copy");
      }

      /* unselect the range */
      testingCodeToCopy.setAttribute("type", "hidden");
      window.getSelection().removeAllRanges();
    },
    getTokens: function() {
      var self = this;
      window.backend.retrieveTokenAmount().then(result => {
        self.tokenAmount = result;
      });
    },
    addNotification(verticalAlign, horizontalAlign, color, copied) {
      setTimeout(() => {
        this.$notifications.clear();
      }, 6000);
      this.$notify({
        component: copied,
        icon: "ti-check",
        horizontalAlign: horizontalAlign,
        verticalAlign: verticalAlign,
        type: this.type[color],
        onClick: () => {
          this.$notifications.clear();
        }
      });
    }
  },
  computed: {
    ...mapState('wallet', ['tokenAmount', 'currency', 'totalValue', 'address']),
    ...mapState('dashboard', ['counters', 'toggle', 'stat', 'chart']) 
  },
  filters: {
    asCurrency: function(value, currency) {

      if (currency == "") return "";

      var formatter
      if (currency == "SHI") {
        formatter = new Intl.NumberFormat(navigator.language);
      } else if (currency == "BTC") {
        formatter = new Intl.NumberFormat(navigator.language, {
          style: "currency",
          currency: "XBT",
          minimumFractionDigits: 2,
          maximumFractionDigits: 8
        });
      } else {
        formatter = new Intl.NumberFormat(navigator.language, {
          style: "currency",
          currency: currency,
          minimumFractionDigits: 2,
          maximumFractionDigits: 2
        });
      }
      return formatter.format(value).replace(/XBT/,'₿');
    }
  },

  /**
   * Chart data used to render stats, charts. Should be replaced with server data
   */

  data() {
    return {
      type: ["", "info", "success", "warning", "danger"],
      notifications: {
        topCenter: false
      },
      transactionChart: {
        options: {
          low: 0,
          high: 1000,
          showArea: true,
          height: "15.3125em",
          axisX: {
            showGrid: false
          },
          lineSmooth: Chartist.Interpolation.simple({
            divisor: 3	
          }),
          showLine: true,
          showPoint: false
        },
      },
      throughputChart: {
        options: {
          seriesBarDistance: 10,
          axisX: {
            showGrid: false
          },
          height: "15.3125em"
        }
      }
    };
  }
};

</script>

<style>
.text-overflow {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.stats-card .card-body .row [class^="col"] {
  margin-left: 0;
  margin-right: 0;
  padding-left: 0;
  padding-right: 0;
}

.stats-card .card-body .row [class^="col"] .numbers {
  margin-top: 0.5rem;
}

.card-footer {
  margin-top: auto;
}

.wallet-address > p-button {
  margin-bottom: 10em;
}
</style>
