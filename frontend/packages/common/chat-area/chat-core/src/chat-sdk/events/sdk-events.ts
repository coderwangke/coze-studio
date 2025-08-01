/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
 
export enum SdkEventsEnum {
  MESSAGE_RECEIVED_AND_UPDATE = 'message_received_and_update',
  /**
   * 监测拉取回复消息状态变化
   */
  MESSAGE_PULLING_STATUS = 'message_pulling_status',
  ERROR = 'error',
}
