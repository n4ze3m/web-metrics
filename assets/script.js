(function(window, document) {
    var Analytics = {
        init: function(options) {
            if (typeof window === 'undefined') return;

            this.options = options || {};
            this.endpoint = options.endpoint || '/events';
            this.trackPageviews = options.trackPageviews !== false;
            this.websiteId = options.websiteId;

            if(!this.websiteId) {
                throw new Error('Website ID is required');
            }

            this.userId = this.getUserId();

            if (this.trackPageviews) {
                this.trackPageview();
            }

            if (options.trackClicks) {
                this.trackClicks();
            }

            if (typeof window.next !== 'undefined') {
                window.next.router.events.on('routeChangeComplete', this.trackPageview.bind(this));
            }

            this.startTimeOnPage = new Date();
            window.addEventListener('beforeunload', this.trackTimeOnPage.bind(this));
        },

        generateUserId: function() {
            return 'user_' + Math.random().toString(36).substr(2, 9);
        },

        getUserId: function() {
            var userId = this.getCookie('_wm_user_id');
            if (!userId) {
                userId = this.generateUserId();
                this.setCookie('_wm_user_id', userId, 365);
            }
            return userId;
        },

        getCookie: function(name) {
            var match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
            return match ? match[2] : null;
        },

        setCookie: function(name, value, days) {
            var expires = '';
            if (days) {
                var date = new Date();
                date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
                expires = '; expires=' + date.toUTCString();
            }
            document.cookie = name + '=' + value + expires + '; path=/';
        },

        trackEvent: function(eventName, eventData) {
            if (typeof window === 'undefined') return Promise.resolve();

            var data = {
                user_id: this.userId,
                website_id: this.websiteId,
                event_type: eventName,
                timestamp: new Date().toISOString(),
                url: window.location.href,
                referrer: document.referrer,
                user_agent: navigator.userAgent,
                screen_width: window.screen.width,
                screen_height: window.screen.height,
                ...eventData
            };

            return this.sendRequest(data);
        },

        trackPageview: function() {
            return this.trackEvent('pageview');
        },

        trackClicks: function() {
            var self = this;
            document.addEventListener('click', function(e) {
                var target = e.target;
                while (target && target.tagName !== 'A') {
                    target = target.parentNode;
                }
                if (target && target.tagName === 'A') {
                    self.trackEvent('click', {
                        data: JSON.stringify({
                            text: target.innerText,
                            href: target.href
                        })
                    });
                }
            });
        },

        trackTimeOnPage: function() {
            var endTime = new Date();
            var timeSpent = endTime - this.startTimeOnPage;
            this.trackEvent('time_on_page', {
                data: JSON.stringify({
                    seconds: Math.floor(timeSpent / 1000)
                })
            });
        },

        sendRequest: function(data) {
            return new Promise(function(resolve, reject) {
                var xhr = new XMLHttpRequest();
                xhr.open('POST', this.endpoint, true);
                xhr.setRequestHeader('Content-Type', 'application/json');
                xhr.onload = function() {
                    if (xhr.status === 200) {
                        resolve(JSON.parse(xhr.responseText));
                    } else {
                        reject(Error(xhr.statusText));
                    }
                };
                xhr.onerror = function() {
                    reject(Error('Network Error'));
                };
                xhr.send(JSON.stringify(data));
            }.bind(this));
        }
    };

    if (typeof window !== 'undefined') {
        window.Analytics = Analytics;
    }
})(typeof window !== 'undefined' ? window : {}, typeof document !== 'undefined' ? document : {});